# CMDR
[![test](https://github.com/MrLYC/cmdr/actions/workflows/unittest.yml/badge.svg)](https://github.com/MrLYC/cmdr/actions/workflows/unittest.yml) [![release](https://github.com/MrLYC/cmdr/actions/workflows/release.yml/badge.svg)](https://github.com/MrLYC/cmdr/actions/workflows/main.yml) [![codecov](https://codecov.io/gh/MrLYC/cmdr/branch/master/graph/badge.svg?token=mo4TJP4mQt)](https://codecov.io/gh/MrLYC/cmdr)

CMDR is a simple command version management tool that helps you quickly switch from multiple command versions.

## Installation
1. Download the latest release from [GitHub](https://github.com/MrLYC/cmdr/releases);
2. Make sure the download asset is executable;
3. Run the following command to install the binary:
    ```shell
    % /path/to/cmdr init
    ```
4. Restart your shell and run the following command to verify the installation:
    ```shell
    % cmdr version
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
To upgrade the CMDR, just run:
```shell
% cmdr upgrade
```