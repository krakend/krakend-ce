#!/bin/sh

OS=$(uname)
GLIBC=UNKNOWN-0.0.0

case $OS in
Linux*)
    if ldd --version 2>&1 | grep -i musl > /dev/null; then
        GLIBC=MUSL-$(ldd --version 2>&1 | grep Version | cut -d" " -f2)
    else
        GLIBC=GLIBC-$(ldd --version 2>&1  | grep ^ldd | awk '{print $(NF)}')
    fi
    ;;
Darwin*)
    GLIBC=DARWIN-$(sw_vers | grep ProductVersion | cut -d$'\t' -f2)
    ;;
*)
  ;;
esac

echo $GLIBC