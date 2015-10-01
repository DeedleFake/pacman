package pacman

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	verRE = regexp.MustCompile(`([0-9]+:)?[a-zA-Z0-9._]+-[0-9]+`)
)

type Version string

func (v Version) Valid() bool {
	return verRE.MatchString(string(v))
}

func (v Version) Epoch() int {
	parts := strings.SplitN(string(v), ":", 2)
	if len(parts) <= 1 {
		return 0
	}

	e, _ := strconv.ParseInt(parts[0], 10, 0)
	return int(e)
}

func (v Version) Pkgver() string {
	parts := strings.SplitN(string(v), ":", 2)
	if len(parts) > 1 {
		parts = parts[1:]
	}

	parts = strings.SplitN(parts[0], "-", 2)

	return parts[0]
}

func (v Version) Pkgrel() int {
	parts := strings.SplitN(string(v), "-", 2)
	if len(parts) <= 1 {
		return 1
	}

	r, _ := strconv.ParseInt(parts[1], 10, 0)
	return int(r)
}
