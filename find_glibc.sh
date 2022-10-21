#!/bin/sh

OS=$(uname)
GLIBC=UNKNOWN-0.0.0

get_os_version() {
    . /etc/os-release
    os_release="${ID}-${VERSION_ID}"
}

case $OS in
Linux*)
    if ldd --version 2>&1 | grep -i musl > /dev/null; then
        get_os_version
        GLIBC="MUSL-$(ldd --version 2>&1 | grep Version | cut -d\  -f2)_($os_release)"
    else
        get_os_version
        GLIBC="GLIBC-$(ldd --version 2>&1  | grep ^ldd | awk '{print $(NF)}')_($os_release)"
    fi
    ;;
Darwin*)
    GLIBC=DARWIN-$(sw_vers | grep ProductVersion | cut -d$'\t' -f2)
    ;;
*)
  ;;
esac

echo $GLIBC