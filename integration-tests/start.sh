#!/usr/bin/env bash

root_dir="$(dirname $0)"

set -e
export CMDR_LOG_LEVEL=debug
export CMDR_CORE_CONFIG_PATH=/tmp/cmdr.yaml
export CMDR_CORE_ROOT_DIR="$(pwd)/.cmdr"

# Put generated profile under CMDR_CORE_ROOT_DIR to avoid modifying repo-tracked files.
if [ -z "${CMDR_CORE_PROFILE_DIR}" ]; then
  export CMDR_CORE_PROFILE_DIR="${CMDR_CORE_ROOT_DIR}/profile"
fi

if [ "${SKIP_INSTALL}" != "true" ]; then
  ./install.sh -d 30
else
  # When skipping install, ensure profile is set up
  if [ ! -f "./cmdr" ]; then
    echo "Error: ./cmdr binary not found. Run 'go build -o cmdr .' first."
    exit 1
  fi
  if [ ! -d "${CMDR_CORE_PROFILE_DIR}" ]; then
    ./cmdr init
  fi
fi

chmod +x \
  "$root_dir/cmd_v1.sh" \
  "$root_dir/cmd_v2.sh" \
  "$root_dir/cmd_v3.sh" \
  "$root_dir/cmd_v4.sh" \
  "$root_dir/cmd_v5.sh" \
  "$root_dir/cmd_v6.sh" 2>/dev/null || true

source "${CMDR_CORE_PROFILE_DIR}/cmdr_initializer.sh"

# Bash can cache command paths. Since cmdr activates by updating symlinks,
# clear the command hash table to make sure we always execute latest links.
hash -r 2>/dev/null || true

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
  hash -r 2>/dev/null || true
  cmd | grep "v$v"
done

cmdr command unset -n cmd
hash -r 2>/dev/null || true
cmd && false || true

rm -rf "./.cmdr/shims/cmd/cmd_1.0.0"
cmdr command list -n cmd -v "1.0.0"

cmdr doctor
cmdr command list -n cmd -v "1.0.0" && false || true

cmdr command list -n cmd -v "2.0.0"
cmdr command list -n cmd -v "3.0.0"
cmdr command use -n cmd -v "2.0.0"
hash -r 2>/dev/null || true
cmd | grep "v2.0.0"

cmdr command unset -n cmd
hash -r 2>/dev/null || true

cmdr command remove -n cmd -v "2.0.0"
cmdr command remove -n cmd -v "3.0.0"

## cmdr clean integration test
cmdr command install -a -n cmdclean -v "1.0.0" -l "$root_dir/cmd_v1.sh"
cmdr command install -n cmdclean -v "2.0.0" -l "$root_dir/cmd_v2.sh"
cmdr command install -n cmdclean -v "3.0.0" -l "$root_dir/cmd_v3.sh"
cmdr command install -n cmdclean -v "4.0.0" -l "$root_dir/cmd_v4.sh"
cmdr command install -n cmdclean -v "5.0.0" -l "$root_dir/cmd_v5.sh"
cmdr command define -a -n cmdclean -v "6.0.0" -l "$root_dir/cmd_v6.sh"

if [ "$(uname -s)" = "Darwin" ]; then
  clean_trash_dir="$HOME/.Trash/cmdr-cleaned"
else
  clean_trash_dir="/tmp/cmdr-cleaned"
fi

rm -rf "$clean_trash_dir/cmdclean" 2>/dev/null || true

python3 - <<'PY'
import os
import time

root = os.environ.get('CMDR_CORE_ROOT_DIR')
if not root:
    raise RuntimeError('CMDR_CORE_ROOT_DIR is not set')

base = os.path.join(root, 'shims', 'cmdclean')
items = [
    ('cmdclean_1.0.0', 300),
    ('cmdclean_2.0.0', 200),
    ('cmdclean_3.0.0', 150),
    ('cmdclean_4.0.0', 120),
    ('cmdclean_5.0.0', 90),
]

now = time.time()
for filename, days in items:
    p = os.path.join(base, filename)
    ts = now - days * 86400
    os.utime(p, (ts, ts))
PY

cmdclean | grep "v6.0.0"

hash -r 2>/dev/null || true

cmdr clean -n cmdclean --age 100 --keep 3

# should fail for unknown command name
if cmdr clean -n cmdr_clean_not_exist --age 100 --keep 3; then
  echo "expected cmdr clean to fail for unknown name"
  exit 1
fi

cmdr command list -n cmdclean -v "1.0.0" 2>&1 && { echo "ERROR: cmdclean v1 should be cleaned"; exit 1; } || echo "✓ cmdclean v1 cleaned"
cmdr command list -n cmdclean -v "2.0.0" 2>&1 && { echo "ERROR: cmdclean v2 should be cleaned"; exit 1; } || echo "✓ cmdclean v2 cleaned"

cmdr command list -n cmdclean -v "3.0.0" > /dev/null || { echo "ERROR: cmdclean v3 should exist"; exit 1; }
cmdr command list -n cmdclean -v "4.0.0" > /dev/null || { echo "ERROR: cmdclean v4 should exist"; exit 1; }
cmdr command list -n cmdclean -v "5.0.0" > /dev/null || { echo "ERROR: cmdclean v5 should exist"; exit 1; }
cmdr command list -n cmdclean -v "6.0.0" > /dev/null || { echo "ERROR: cmdclean v6 should exist"; exit 1; }

test -f "$clean_trash_dir/cmdclean/cmdclean_1.0.0" || { echo "ERROR: missing trashed shim for cmdclean v1"; exit 1; }
test -f "$clean_trash_dir/cmdclean/cmdclean_2.0.0" || { echo "ERROR: missing trashed shim for cmdclean v2"; exit 1; }

cmdclean | grep "v6.0.0"

echo "All tests passed successfully!"
