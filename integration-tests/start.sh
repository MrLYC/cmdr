#!/usr/bin/env bash

root_dir="$(dirname $0)"

set -e

./install.sh

export CMDR_LOG_LEVEL=debug CMDR_CORE_CONFIG_PATH=/tmp/cmdr.yaml
cmdr config set -k core.root_dir -v "$(pwd)/.cmdr"
cmdr config set -k core.profile_dir -v "$(pwd)/profile"

cmdr init
source ./profile/cmdr_initializer.sh

set -x

cmdr command list -n cmdr
cmdr config list

cmdr command install -a -n cmd -v "1.0.0" -l "$root_dir/cmd_v1.sh"
cmdr command install -n cmd -v "2.0.0" -l "$root_dir/cmd_v2.sh"
cmdr command define -n cmd -v "3.0.0" -l "$root_dir/cmd_v3.sh"

cmd | grep "v1.0.0"

for v in "2.0.0" "3.0.0"; do
  cmdr command use -n cmd -v "$v"
  cmd | grep "v$v"
done

cmdr command unset -n cmd
cmd && false || true

rm -rf "./.cmdr/shims/cmd/cmd_1.0.0"
cmdr command list -n cmd -v "1.0.0"
cmdr doctor
cmdr command list -n cmd -v "1.0.0" && false || true

cmdr command remove -n cmd -v "2.0.0"
cmdr command remove -n cmd -v "3.0.0"

current_version=$(cmdr version)

cmdr upgrade
newest_version=$(cmdr version)

# make sure cmdr has been upgraded
test "${current_version}" != "${newest_version}"
cmdr command list -n cmdr -a -v "${newest_version}"

branch="${GITHUB_BRANCH:-${GITHUB_REF##*/}}"
cmdr command install -a -n cmdr -v "0.0.0" -l "go://github.com/MrLYC/cmdr@${branch}"
cmdr init

cmdr init --upgrade
cmdr command list -n cmdr

activated_version=$(cmdr version)

test "${current_version}" == "${activated_version}"
cmdr command list -n cmdr -a -v "${current_version}"