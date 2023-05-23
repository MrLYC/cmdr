#!/usr/bin/env bash

set -e

max_retry=3
retry_delay=1
tag_name=""
ghproxy=""

while getopts "r:d:t:p" opt; do
    case $opt in
        r)  max_retry="${OPTARG}"  ;;
        d)  retry_delay="${OPTARG}" ;;
        t)  tag_name="${OPTARG}"   ;;
        p)  ghproxy="1"    ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            exit 1
            ;;
    esac
done

echo "${ghproxy}"

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

function update_tag_name() {
    semver='v\d+\.\d+\.\d+'

    tag_name=$(curl -s https://api.github.com/repos/MrLYC/cmdr/releases/latest | grep 'tag_name' | grep -o "v[^"]*")
    if [[ -n "${tag_name}" ]]; then
        return 0
    fi

    tag_name=$(curl -s https://github.com/MrLYC/cmdr/releases.atom | grep '<title>' | grep -o 'v[^<]*' -m 1)
    if [[ -n "${tag_name}" ]]; then
        return 0
    fi

    return 1
}

function download_cmdr() {
    if [[ -z "${tag_name}" ]]; then
        echo "Quering cmdr latest release for ${os}/${arch}..."
        update_tag_name
    fi

    if [[ -z "${tag_name}" ]]; then
        echo "Failed to query cmdr latest release for ${os}/${arch}"
        return 1
    fi

    download_url="https://github.com/MrLYC/cmdr/releases/download/${tag_name}/cmdr_${goos}_${goarch}"

    if [[ "${ghproxy}" == "1" ]]; then
        download_url="https://ghproxy.com/${download_url}"
    fi

    target="$1"
    echo "Downloading cmdr ${tag_name}..."
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
    elif [[ "${retry}" = "${max_retry}" ]]; then
        echo "Activating proxy automaticly..."
        ghproxy="1"
    fi

    echo "Failed to download cmdr, retry after ${retry_delay} second"
    retry=$((retry+1))
    sleep "${retry_delay}"
done

echo "Initializing cmdr..."
"${target}" init
"${target}" command list -n cmdr

rm -f "${target}"
echo "Please restart your terminal to activate the cmdr command"