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

function download_cmdr() {
    if [[ -z "${tag_name}" ]]; then
        echo "Quering cmdr latest release for ${os}/${arch}..."
        download_url=$(curl -s https://api.github.com/repos/MrLYC/cmdr/releases/latest | grep browser_download_url | grep -o "https://.*/cmdr_${goos}_${goarch}"
    else
        download_url="https://github.com/MrLYC/cmdr/releases/download/${tag_name}/cmdr_${goos}_${goarch}"
    fi

    if [[ -z "${download_url}" ]]; then
        echo "Failed to get cmdr download url"
        return 1
    fi

    target="$1"
    echo "Downloading cmdr (${download_url})..."
    curl -L -o "${target}" "${download_url}"
    chmod +x "${target}"
}

echo "Downloading cmdr..."
target="/tmp/cmdr_${RANDOM}"
retry=0
until download_cmdr "${target}"; do
    if [[ "${retry}" -gt "${max_retry}" ]]; then
        echo "Failed to download cmdr"
        exit 1
    fi

    echo "Failed to download cmdr, retry after 1 second"
    retry=$((retry+1))
    sleep 1
done

echo "Initializing cmdr..."
"${target}" init
"${target}" command list -n cmdr

rm -f "${target}"
echo "Please restart your terminal to activate the cmdr command"