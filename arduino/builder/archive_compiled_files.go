// This file is part of arduino-cli.
//
// Copyright 2023 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package builder

import (
	"github.com/arduino/go-paths-helper"
	"github.com/pkg/errors"
)

// ArchiveCompiledFiles fixdoc
func (b *Builder) archiveCompiledFiles(buildPath *paths.Path, archiveFile *paths.Path, objectFilesToArchive paths.PathList) (*paths.Path, error) {
	archiveFilePath := buildPath.JoinPath(archiveFile)

	if b.onlyUpdateCompilationDatabase {
		if b.logger.Verbose() {
			b.logger.Info(tr("Skipping archive creation of: %[1]s", archiveFilePath))
		}
		return archiveFilePath, nil
	}

	if archiveFileStat, err := archiveFilePath.Stat(); err == nil {
		rebuildArchive := false
		for _, objectFile := range objectFilesToArchive {
			objectFileStat, err := objectFile.Stat()
			if err != nil || objectFileStat.ModTime().After(archiveFileStat.ModTime()) {
				// need to rebuild the archive
				rebuildArchive = true
				break
			}
		}

		// something changed, rebuild the core archive
		if rebuildArchive {
			if err := archiveFilePath.Remove(); err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			if b.logger.Verbose() {
				b.logger.Info(tr("Using previously compiled file: %[1]s", archiveFilePath))
			}
			return archiveFilePath, nil
		}
	}

	for _, objectFile := range objectFilesToArchive {
		properties := b.buildProperties.Clone()
		properties.Set("archive_file", archiveFilePath.Base())
		properties.SetPath("archive_file_path", archiveFilePath)
		properties.SetPath("object_file", objectFile)

		command, err := b.prepareCommandForRecipe(properties, "recipe.ar.pattern", false)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if err := b.execCommand(command); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return archiveFilePath, nil
}
