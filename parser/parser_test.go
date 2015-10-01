package parser_test

import (
	"."
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	const file = `%NAME% This shouldn't be here.
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

	ex := []parser.Token{
		parser.Header("NAME"),
		parser.Entry("pacman"),
		parser.Header("VERSION"),
		parser.Entry("4.2.1-2"),
		parser.Header("DESC"),
		parser.Entry("A library-based package manager with dependency support"),
		parser.Header("URL"),
		parser.Entry("http://www.archlinux.org/pacman/"),
		parser.Header("ARCH"),
		parser.Entry("x86_64"),
		parser.Header("BUILDDATE"),
		parser.Entry("1437702809"),
		parser.Header("INSTALLDATE"),
		parser.Entry("1438007203"),
		parser.Header("PACKAGER"),
		parser.Entry("Allan McRae <allan@archlinux.org>"),
		parser.Header("SIZE"),
		parser.Entry("4579328"),
		parser.Header("GROUPS"),
		parser.Entry("base"),
		parser.Entry("base-devel"),
		parser.Header("LICENSE"),
		parser.Entry("GPL"),
		parser.Header("VALIDATION"),
		parser.Entry("pgp"),
		parser.Header("REPLACES"),
		parser.Entry("pacman-contrib"),
		parser.Header("DEPENDS"),
		parser.Entry("bash"),
		parser.Entry("glibc"),
		parser.Entry("libarchive>=3.1.2"),
		parser.Entry("curl>=7.39.0"),
		parser.Entry("gpgme"),
		parser.Entry("pacman-mirrorlist"),
		parser.Entry("archlinux-keyring"),
		parser.Header("CONFLICTS"),
		parser.Entry("pacman-contrib"),
		parser.Header("PROVIDES"),
		parser.Entry("pacman-contrib"),
	}

	p := parser.New(strings.NewReader(file))
	for i := 0; p.Next() && i < len(ex); i++ {
		if p.Tok() != ex[i] {
			t.Errorf("%v: Expected %q. Got %q.", i, ex[i], p.Tok())
		}
	}
	if err := p.Err(); err != nil {
		t.Fatal(err)
	}
}
