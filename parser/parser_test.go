package parser_test

import (
	"."
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	const file = `%NAME%
pacman

%VERSION%
4.2.1-2

%DESC%
A library-based package manager with dependency support

%URL%
http://www.archlinux.org/pacman/

%ARCH%
x86_64

%BUILDDATE%
1437702809

%INSTALLDATE%
1438007203

%PACKAGER%
Allan McRae <allan@archlinux.org>

%SIZE%
4579328

%GROUPS%
base
base-devel

%LICENSE%
GPL

%VALIDATION%
pgp

%REPLACES%
pacman-contrib

%DEPENDS%
bash
glibc
libarchive>=3.1.2
curl>=7.39.0
gpgme
pacman-mirrorlist
archlinux-keyring

%CONFLICTS%
pacman-contrib

%PROVIDES%
pacman-contrib`

	p := parser.New(strings.NewReader(file))
	for p.Next() {
		t.Log(p.Tok())
	}
	if err := p.Err(); err != nil {
		t.Fatal(err)
	}
}
