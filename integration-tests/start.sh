#!/usr/bin/env bash

root_dir="$(dirname $0)"

set -e

go build -o cmdr .

./cmdr init
source ~/.cmdr/profile/cmdr_initializer.sh

set -x

cmdr version
cmdr config list

cmdr command install -a -n cmd -v "1.0.0" -l "$root_dir/cmd_v1.sh"
cmdr command install -n cmd -v "2.0.0" -l "$root_dir/cmd_v2.sh"
cmdr command define -n cmd -v "3.0.0" -l "$root_dir/cmd_v3.sh"

cmd | grep "v1"

for v in "2.0.0" "3.0.0"; do
  cmdr command use -n cmd -v "$v"
  cmd | grep "$v"
done

cmdr command unset -n cmd
cmd && false || true

cmdr command uninstall -n cmd -v "1.0.0"
cmdr command uninstall -n cmd -v "2.0.0"
cmdr command undefine -n cmd -v "3.0.0"

cmdr upgrade