# CMDR
[![test](https://github.com/MrLYC/cmdr/actions/workflows/test.yml/badge.svg)](https://github.com/MrLYC/cmdr/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/MrLYC/cmdr/branch/master/graph/badge.svg?token=mo4TJP4mQt)](https://codecov.io/gh/MrLYC/cmdr) ![Go version](https://img.shields.io/github/go-mod/go-version/mrlyc/cmdr) ![release](https://img.shields.io/github/v/release/mrlyc/cmdr?label=version)

CMDR is a simple command version management tool that helps you quickly switch from multiple command versions.

## Installation

### Script
Run one for the following command to install the latest version of CMDR:

```shell
curl -o- https://raw.githubusercontent.com/MrLYC/cmdr/master/install.sh | ${SHELL:-bash}
```

### Manual
1. Download the latest release from [GitHub](https://github.com/MrLYC/cmdr/releases/latest);
2. Make sure the download asset is executable;
3. Run the following command to install the binary:
    ```shell
    /path/to/cmdr init
    ```
4. Restart your shell and run the following command to verify the installation:
    ```shell
    cmdr version
    ```

## Get Started
To install a new command, run the following command:
```shell
cmdr command install -n <command-name> -v <version> -l <path-or-url>
```

Then you can list the installed commands by running the following command:
```shell
cmdr command list -n <command-name>
```

Use a specified command version:
```shell
cmdr command use -n <command-name> -v <version>
```

## Upgrade
To upgrade the CMDR, just run:
```shell
cmdr upgrade
```