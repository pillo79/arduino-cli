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
	"testing"

	semver "go.bug.st/relaxed-semver"

	"github.com/stretchr/testify/require"
)

func TestLibrarySlug(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"PlainName", "Bridge", "BRIDGE"},
		{"AlreadyUpper", "WIFI101", "WIFI101"},
		{"SpacesAndDash", "Bridge Client-2", "BRIDGE_CLIENT_2"},
		{"Dots", "Adafruit.SSD1306", "ADAFRUIT_SSD1306"},
		{"ConsecutiveSeparators", "Foo   Bar--Baz", "FOO_BAR_BAZ"},
		{"LeadingDigit", "101Lib", "_101LIB"},
		{"Empty", "", ""},
		{"OnlySeparators", "---", "_"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, librarySlug(c.in))
		})
	}
}

func TestEncodeLibraryVersion(t *testing.T) {
	parse := func(s string) *semver.Version {
		v, err := semver.Parse(s)
		require.NoError(t, err)
		return v
	}

	cases := []struct {
		name string
		in   *semver.Version
		want uint32
	}{
		{"Full", parse("1.2.3"), 0x010203},
		{"MajorMinorOnly", parse("1.2"), 0x010200},
		{"MajorOnly", parse("1"), 0x010000},
		{"Nil", nil, 1},
		{"EmptyRawVersion", parse(""), 1},
		{"OverflowComponentClamped", parse("300.0.0"), 0xFF0000},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, encodeLibraryVersion(c.in))
		})
	}
}
