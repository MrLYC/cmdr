#!/usr/bin/env bash

set -e
max_retry=3
tag_name=""

while getopts "r:t:" opt; do
    case $opt in
        r)  max_retry="${OPTARG}"  ;;
        t)  tag_name="${OPTARG}"   ;;
    esac
done

shift $((OPTIND-1))

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

if [[ -z "${tag_name}" ]]; then
    echo "Quering cmdr latest release for ${os}/${arch}..."
    tag_name=$(curl --silent https://api.github.com/repos/MrLYC/cmdr/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
fi

target="/tmp/cmdr_${RANDOM}"
download_url="https://github.com/MrLYC/cmdr/releases/download/${tag_name}/cmdr_${goos}_${goarch}"
echo "Downloading cmdr ${tag_name} (${download_url})..."
retry=0
until curl -L -o "${target}" "${download_url}"; do
    echo "Retry after 1 second..."
    sleep 1
    retry=$((retry + 1))
    if [[ ${retry} -gt "${max_retry}" ]]; then
        echo "Failed to download cmdr ${tag_name}"
        exit 1
    fi
done
chmod +x "${target}"

echo "Initializing cmdr..."
"${target}" init
"${target}" command list -n cmdr

rm -f "${target}"
echo "Please restart your terminal to activate the cmdr command"