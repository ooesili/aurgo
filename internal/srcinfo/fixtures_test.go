package srcinfo_test

var fixtureSinglePackge = `
# super important comment
pkgbase = dopepkg
	arch = x86_64
	depends = libdope
	depends = leftpad
	makedepends = cmake
	makedepends = maven
	checkdepends = checktool
	checkdepends = testlib

pkgname = dopepkg
`

var fixtureMultiArch = `
pkgbase = multiarchpkg
	arch = x86_64
	arch = i686

pkgname = multiarchpkg
`

var fixtureAnyArch = `
pkgbase = anyarchpkg
	arch = any

pkgname = anyarchpkg
`

var fixtureArchSpecificDeps = `
pkgbase = archspecificdeps
	arch = x86_64
	arch = i686
	depends = leftpad
	depends_x86_64 = libnice64
	makedepends = cmake
	makedepends_x86_64 = maven
	checkdepends = checktool
	checkdepends_x86_64 = testlib

pkgname = archspecificdeps
`

var fixtureSplitPackage = `
pkgbase = splitpkg
	arch = any
	depends = libdope
	makedepends = cmake
	checkdepends = testlib

pkgname = splitpkg-gtk
	depends = leftpad
	makedepends = maven
	checkdepends = checktool

pkgname = libsplit
	depends = glib2
	makedepends = tup
	checkdepends = check
`

var fixtureVersionConstraints = `
pkgbase = versiondeps
	arch = x86_64
	depends = dep1>0.5.1
	depends = dep2<0.5.1
	depends = dep3>=0.5.1
	depends = dep4<=0.5.1
	depends = dep5=0.5.1
	makedepends = makedep1>0.5.1
	makedepends = makedep2<0.5.1
	makedepends = makedep3>=0.5.1
	makedepends = makedep4<=0.5.1
	makedepends = makedep5=0.5.1
	checkdepends = checkdep1>0.5.1
	checkdepends = checkdep2<0.5.1
	checkdepends = checkdep3>=0.5.1
	checkdepends = checkdep4<=0.5.1
	checkdepends = checkdep5=0.5.1

pkgname = versiondeps
`
