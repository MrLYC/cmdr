#!/usr/bin/env bash

case "$(uname -s)" in
    Linux)
        goos=linux
        ;;
    Darwin)
        goos=darwin
        ;;
    *)
        echo "Unsupported OS $(uname -a)"
        exit 1
        ;;
esac

case "$(uname -m)" in
    x86_64 | x64 | amd64)
        goarch=amd64
        ;;
    arm64 | aarch64 | armv8*)
        goarch=arm64
        ;;
    x86 | i686 | i386)
        goarch=386
        ;;
    armv5*)
        goarch=armv5
        ;;
    armv6*)
        goarch=armv6
        ;;
    armv7*)
        goarch=armv7
        ;;
    *)
        echo "Unsupported ARCH $(uname -a)"
        exit 1
        ;;
esac

set -ex
target="/tmp/cmdr_${RANDOM}"
download_url=$(curl https://api.github.com/repos/MrLYC/cmdr/releases/latest 2>/dev/null | grep browser_download_url |  grep -o "https://.*/cmdr_${goos}_${goarch}")
curl -L -o "${target}" "${download_url}"
chmod +x "${target}"

"${target}" init
"${target}" command list -n cmdr

rm -f "${target}"
set +x

echo "restart your terminal to activate the cmdr command"