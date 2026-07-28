// This file is part of arduino-cli.
//
// Copyright (C) Arduino s.r.l. and/or its affiliated companies
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html

package builder

import (
	"regexp"
	"strings"
)

var (
	invalidIdentifierChars = regexp.MustCompile(`[^A-Za-z0-9_]+`)
	leadingDigit           = regexp.MustCompile(`^[0-9]`)
)

// librarySlug converts a library name into a valid uppercase C identifier.
// Runs of characters outside [A-Za-z0-9_] collapse to a single "_"; the
// result is uppercased, and a leading "_" is added if it would otherwise
// start with a digit.
func librarySlug(name string) string {
	slug := invalidIdentifierChars.ReplaceAllString(name, "_")
	slug = strings.ToUpper(slug)
	if leadingDigit.MatchString(slug) {
		slug = "_" + slug
	}
	return slug
}
