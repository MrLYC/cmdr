#!/usr/bin/env bash

root_dir="$(dirname $0)"
export CMDR_LOG_LEVEL=debug CMDR_CORE_CONFIG_PATH=/tmp/cmdr.yaml CMDR_CORE_ROOT_DIR=$(pwd)/.cmdr CMDR_CORE_PROFILE_DIR=$(pwd)/profile

set -e

./install.sh -d 30

source ./profile/cmdr_initializer.sh

set -x

newest_version=$(cmdr version)

commit_hash="$(git rev-parse HEAD)"
cmdr command install -a -n cmdr -v "0.0.0" -l "go://github.com/mrlyc/cmdr@${commit_hash}"
cmdr init --upgrade

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

cmdr command remove -n cmdr -v "${newest_version}"

cmdr upgrade
activated_version=$(cmdr version)

# make sure cmdr has been upgraded

test "${newest_version}" == "${activated_version}"