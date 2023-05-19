#!/usr/bin/env bash

os="$(uname -s)"
case "${os}" in
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

arch="$(uname -m)"
case "${arch}" in
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

echo "Downloading cmdr ${os}/${arch}..."

set -ex
target="/tmp/cmdr_${RANDOM}"
tag_name=$(curl --silent https://api.github.com/repos/MrLYC/cmdr/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -L -o "${target}" "https://github.com/MrLYC/cmdr/releases/download/${tag_name}/cmdr_${goos}_${goarch}"
chmod +x "${target}"

"${target}" init
"${target}" command list -n cmdr

rm -f "${target}"
set +x

echo "restart your terminal to activate the cmdr command"