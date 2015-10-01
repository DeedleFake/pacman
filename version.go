package pacman

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	verRE = regexp.MustCompile(`([0-9]+:)?[a-zA-Z0-9._]+-[0-9]+`)
)

// Version is a Pacman package version string. For more information,
// see https://wiki.archlinux.org/index.php/PKGBUILD#Version
type Version string

// Valid checks if the version string is valid.
func (v Version) Valid() bool {
	return verRE.MatchString(string(v))
}

// Epoch returns the epoch part of the package version. If the version
// doesn't contain an epoch, it returns -1.
//
// Note: This method doesn't perform any checks. Call v.Valid() before
// calling this if you're unsure if the version string is valid.
func (v Version) Epoch() int {
	parts := strings.SplitN(string(v), ":", 2)
	if len(parts) <= 1 {
		return -1
	}

	e, _ := strconv.ParseInt(parts[0], 10, 0)
	return int(e)
}

// Pkgver returns the main version part of the version string.
//
// Note: This method doesn't perform any checks. Call v.Valid() before
// calling this if you're unsure if the version string is valid.
func (v Version) Pkgver() string {
	parts := strings.SplitN(string(v), ":", 2)
	if len(parts) > 1 {
		parts = parts[1:]
	}

	parts = strings.SplitN(parts[0], "-", 2)

	return parts[0]
}

// Pkgrel returns the package release number. If the version string
// doesn't contain a release number, which is invalid, it returns -1.
//
// Note: This method doesn't perform any checks. Call v.Valid() before
// calling this if you're unsure if the version string is valid.
func (v Version) Pkgrel() int {
	parts := strings.SplitN(string(v), "-", 2)
	if len(parts) <= 1 {
		return -1
	}

	r, _ := strconv.ParseInt(parts[1], 10, 0)
	return int(r)
}
