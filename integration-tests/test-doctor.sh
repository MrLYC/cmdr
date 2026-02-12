#!/usr/bin/env bash

set -e
set -x

root_dir="$(dirname $0)"

export CMDR_LOG_LEVEL=info 
export CMDR_CORE_CONFIG_PATH=/tmp/cmdr-doctor-test.yaml 
export CMDR_CORE_ROOT_DIR=$(pwd)/.cmdr-doctor-test 
export CMDR_CORE_PROFILE_DIR=$(pwd)/profile-doctor-test

cleanup() {
    rm -rf ./.cmdr-doctor-test ./profile-doctor-test /tmp/cmdr-doctor-test.yaml /tmp/.cmdr.backup.* 2>/dev/null || true
}

trap cleanup EXIT

cleanup

if [ ! -f "./cmdr" ]; then
    echo "Building cmdr..."
    go build -o cmdr .
fi

echo "Initializing cmdr..."
./cmdr init

source ./profile-doctor-test/cmdr_initializer.sh

echo "Installing test commands..."
cmdr install -a -n cmd -v "1.0.0" -l "$root_dir/cmd_v1.sh"
cmdr install -n cmd -v "2.0.0" -l "$root_dir/cmd_v2.sh"
cmdr define -n cmd -v "3.0.0" -l "$root_dir/cmd_v3.sh"

echo "Verifying v1 works..."
cmd | grep "v1.0.0"

echo "Simulating missing command by removing v1 shim..."
rm -rf "./.cmdr-doctor-test/shims/cmd/cmd_1.0.0"

echo "Running cmdr doctor..."
cmdr doctor

echo "Verifying v1 was removed..."
cmdr list -n cmd -v "1.0.0" 2>&1 && { echo "ERROR: v1 should not exist"; exit 1; } || echo "✓ v1 correctly removed"

echo "Verifying v2 still exists..."
cmdr list -n cmd -v "2.0.0" > /dev/null || { echo "ERROR: v2 should exist"; exit 1; }
echo "✓ v2 still exists"

echo "Verifying v3 still exists..."
cmdr list -n cmd -v "3.0.0" > /dev/null || { echo "ERROR: v3 should exist"; exit 1; }
echo "✓ v3 still exists"

echo "Verifying v2 can be activated and used..."
cmdr use -n cmd -v "2.0.0"
cmd | grep "v2.0.0" || { echo "ERROR: v2 not working"; exit 1; }
echo "✓ v2 works"

echo "Verifying backup was created in /tmp..."
backup_count=$(ls /tmp/.cmdr*.backup.* 2>/dev/null | wc -l)
if [ "$backup_count" -eq 0 ]; then
    echo "ERROR: No backup found in /tmp"
    exit 1
fi
echo "✓ Backup created in /tmp"

echo ""
echo "=========================================="
echo "All doctor tests passed successfully!"
echo "=========================================="
