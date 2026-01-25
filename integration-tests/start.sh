#!/usr/bin/env bash

root_dir="$(dirname $0)"

set -e
export CMDR_LOG_LEVEL=debug CMDR_CORE_CONFIG_PATH=/tmp/cmdr.yaml CMDR_CORE_ROOT_DIR=$(pwd)/.cmdr CMDR_CORE_PROFILE_DIR=$(pwd)/profile

if [ "${SKIP_INSTALL}" != "true" ]; then
  ./install.sh -d 30
else
  # When skipping install, ensure profile is set up
  if [ ! -f "./cmdr" ]; then
    echo "Error: ./cmdr binary not found. Run 'go build -o cmdr .' first."
    exit 1
  fi
  if [ ! -d "./profile" ]; then
    ./cmdr init
  fi
fi

source ./profile/cmdr_initializer.sh

set -x

newest_version=$(cmdr version)

branch="$(git rev-parse --abbrev-ref HEAD)"
cmdr command install -a -n cmdr -v "0.0.0" -l "./cmdr"
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

cmdr command list -n cmd -v "2.0.0"
cmdr command list -n cmd -v "3.0.0"
cmdr command use -n cmd -v "2.0.0"
cmd | grep "v2.0.0"

cmdr command remove -n cmd -v "2.0.0"
cmdr command remove -n cmd -v "3.0.0"

echo "All tests passed successfully!"