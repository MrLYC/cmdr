#!/usr/bin/env bash

root_dir="$(dirname $0)"

set -e

go build -o cmdr .

export CMDR_LOG_LEVEL=debug CMDR_CORE_CONFIG_PATH=/tmp/cmdr.yaml
./cmdr config set -k core.root_dir -v "$(pwd)/.cmdr"
./cmdr config set -k core.profile_dir -v "$(pwd)/profile"

./cmdr init
source ./profile/cmdr_initializer.sh

set -x

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

version=$(cmdr version)

cmdr upgrade

cmdr version | grep -v "${version}"