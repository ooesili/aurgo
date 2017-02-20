package pacman_test

var fixtureRealOutput = `
Repository      : core
Name            : cronie
Version         : 1.5.1-1
Description     : Daemon that runs specified programs at scheduled times and related tools
Architecture    : x86_64
URL             : https://fedorahosted.org/cronie/
Licenses        : custom:BSD
Groups          : None
Provides        : cron
Depends On      : pam  bash  run-parts
Optional Deps   : pm-utils: defer anacron on battery power
                  smtp-server: send job output via email
                  smtp-forwarder: forward job output to email server
Conflicts With  : cron
Replaces        : None
Download Size   : 68.39 KiB
Installed Size  : 269.00 KiB
Packager        : Gaetan Bisson <bisson@archlinux.org>
Build Date      : Wed 29 Jun 2016 07:26:47 PM CDT
Validated By    : MD5 Sum  SHA-256 Sum  Signature

Repository      : core
Name            : curl
Version         : 7.52.1-2
Description     : An URL retrieval utility and library
Architecture    : x86_64
URL             : https://curl.haxx.se
Licenses        : MIT
Groups          : None
Provides        : libcurl.so=4-64
Depends On      : ca-certificates  krb5  libssh2  openssl  zlib  libpsl
Optional Deps   : None
Conflicts With  : None
Replaces        : None
Download Size   : 850.50 KiB
Installed Size  : 1372.00 KiB
Packager        : Christian Hesse <arch@eworm.de>
Build Date      : Wed 28 Dec 2016 01:16:44 AM CST
Validated By    : MD5 Sum  SHA-256 Sum  Signature

Repository      : core
Name            : grub
Version         : 1:2.02.rc1-1
Description     : GNU GRand Unified Bootloader (2)
Architecture    : x86_64
URL             : https://www.gnu.org/software/grub/
Licenses        : GPL3
Groups          : None
Provides        : grub-common  grub-bios  grub-emu  grub-efi-x86_64
Depends On      : sh  xz  gettext  device-mapper
Optional Deps   : freetype2: For grub-mkfont usage
                  fuse2: For grub-mount usage
                  dosfstools: For grub-mkrescue FAT FS and EFI support
                  efibootmgr: For grub-install EFI support
                  libisoburn: Provides xorriso for generating grub rescue iso using grub-mkrescue
                  os-prober: To detect other OSes when generating grub.cfg in BIOS systems
                  mtools: For grub-mkrescue FAT FS support
Conflicts With  : grub-common  grub-bios  grub-emu  grub-efi-x86_64  grub-legacy
Replaces        : grub-common  grub-bios  grub-emu  grub-efi-x86_64
Download Size   : 5.84 MiB
Installed Size  : 28.79 MiB
Packager        : Christian Hesse <arch@eworm.de>
Build Date      : Fri 03 Feb 2017 01:46:41 PM CST
Validated By    : MD5 Sum  SHA-256 Sum  Signature

Repository      : community
Name            : rust
Version         : 1:1.14.0-1
Description     : Systems programming language focused on safety, speed and concurrency
Architecture    : x86_64
URL             : https://www.rust-lang.org/
Licenses        : MIT  Apache
Groups          : rust
Provides        : None
Depends On      : gcc-libs  llvm-libs
Optional Deps   : None
Conflicts With  : None
Replaces        : None
Download Size   : 33.71 MiB
Installed Size  : 127.25 MiB
Packager        : Alexander Rødseth <rodseth@gmail.com>
Build Date      : Fri 23 Dec 2016 02:06:03 PM CST
Validated By    : MD5 Sum  SHA-256 Sum  Signature
`
