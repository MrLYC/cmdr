#!/usr/bin/env zsh

set -e

/tmp/cmdr init
source ~/.zshrc

cmdr version
cmdr config list

cmdr command install -a -n xsh -v "1.0.0" -l $(which bash)
cmdr command define -n xsh -v "2.0.0" -l $(which zsh)

xsh --version | grep bash

cmdr command use -n xsh -v "2.0.0"

xsh --version | grep zsh

cmdr command unset -n xsh

xsh && false || true

cmdr upgrade