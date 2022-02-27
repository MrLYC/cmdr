# CMDR
[![test](https://github.com/MrLYC/cmdr/actions/workflows/unittest.yml/badge.svg)](https://github.com/MrLYC/cmdr/actions/workflows/unittest.yml) [![release](https://github.com/MrLYC/cmdr/actions/workflows/release.yml/badge.svg)](https://github.com/MrLYC/cmdr/actions/workflows/main.yml) [![codecov](https://codecov.io/gh/MrLYC/cmdr/branch/main/graph/badge.svg?token=mo4TJP4mQt)](https://codecov.io/gh/MrLYC/cmdr)

CMDR is a simple command version management tool that helps you quickly switch from multiple command versions.

## Installation
Download the latest release from [GitHub](https://github.com/MrLYC/cmdr/releases) and make sure it is executable.
Run the following command to install it in your system:
```shell
% /path/to/cmdr setup
```

Check the CMDR version information by running the following command:
```shell
% cmdr version -a
```

## Get Started
To install a new command, run the following command:
```shell
% cmdr command install -n <command-name> -v <version> -l <path_or_url>
```

Then you can list the installed commands by running the following command:
```shell
% cmdr command list -n <command-name>
```

Use a specified command version:
```shell
% cmdr command use -n <command-name> -v <version>
```

## Upgrade
To upgrade the CMDR, run:
```shell
% cmdr upgrade
```