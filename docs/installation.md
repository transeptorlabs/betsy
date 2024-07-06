# Installation

**Requirements to run Betsy**:
1. [Go - >= v1.22.4](https://go.dev/doc/install)
2. [Docker](https://docs.docker.com/engine/install)

## Versioning

Betsy follows [Semantic Versioning](https://semver.org/) for versioning releases. Each release can be found on the repository as a branch with the version number `release/x.y.z.` and a release tag with the version number `vx.y.z`.

The command `betsy --version` allows you to check which version of Betsy you are running. The version number along with the commit hash and commit date should display.

The unstable version(development) format looks as follows:
```shell    
betsy version x.y.z-unstable (abcabcabcabc yyyy-mm-dd)
```

The stable version format looks as follows:
```shell
betsy version x.y.z-stable (abcabcabcabc yyyy-mm-dd)
```

- `Unstable` is for those who want to run the latest version of Betsy to test new features and bug fixes.
- `Stable` is for those who want to run the latest stable version of Betsy.

## Build from the source

> ℹ️ **Info**: Betsy will default to unstable builds(development). If you want to use a stable build, it is recommended to build from a specific release version, `vx.y.z`.

### Linux and Mac

You can clone the [Betsy](https://github.com/transeptorlabs/betsy) repository for UNIX-like operating systems and create a **temporary** build using the command `make betsy`. This method of building requires Go to be installed on your system.

Running `make betsy` creates a standalone executable file in the `betsy/bin` directory. You can run the executable file using the command `./bin/betsy --help`. Or you can move the executable file and run it from another directory.

#### Unstable(development) version build:
```shell
git clone https://github.com/transeptorlabs/betsy.git
cd betsy
make betsy
```

To update the latest development version of Betsy, you can:
1. Stop the CLI(If it is running)
2. Navigate to the Betsy directory 
3. Pull the latest version of the source code from Betsy's Github repository 
4. Build and restart the CLI
   
```shell
cd betsy
git pull
make betsy
```

#### Stable version build:
```shell
git clone https://github.com/transeptorlabs/betsy.git
cd betsy
git checkout vx.y.z
make betsy
```

