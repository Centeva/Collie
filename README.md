# Collie

Collie is a devops tool to help with Unitely.

## Table of Contents

1. [About the Project](#about-the-project)
    * [Built With](#built-with)
2. [Getting Started](#getting-started)
    * [Prerequisites](#prerequisites)
    * [Installation](#installation)
    * [IDE](#ide)
3. [Usage](#usage)
    * [Building local](#building-local)
    * [Docker local](#docker-local)
    * [Private Docker](#private-Docker)
4. [Running Tests](#running-tests)
5. [Deployment](#deployment)
    * [Versioning](#versioning)
6. [Contributing](#contributing)
7. [Related Projects](#related-projects)
8. [Resources](#resources)

## About The Project
This app is built in Go to help with some of the random devops tasks we need to do for Unitely.

Some things this app does:

- format branch names into a k8s compatible format
- connect to a postgresql server and clean up dead databases


### Built With

* [golang](https://golang.org/)
* [docker](https://www.docker.com/)

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

You will need go installed. Either download from [golang](https://golang.org/dl/) or use winget.

```ps
winget install GoLang.Go -v 1.16.6
```

### Installation

1. Clone the repo

    ```sh
    git clone https://bitbucket.org/centeva/Collie.git
    ```

2. Install go modules (not required but will cache modules)

    ```ps
    go mod download
    ```

### IDE

Go does not have an IDE. Most people either use [VSCode](https://code.visualstudio.com/) or Vim. JetBrains also offers [GoLand](https://www.jetbrains.com/go/).

For VSCode there is an extension called `golang.go` that you need to install. Also run `>go.tools.install` and select all 10 of the tools in the list. These are tools the extension uses to check lint, build, etc. By installing them all now we avoid annoying popup notifications.

## Usage

You can this tool with docker or manually.

### Building local
Running `go build` will create a `collie.exe` that you can then run manually. This will work locally but this exe is not cross platform.

### Docker local
You can build the dockerfile locally with `docker build . -t collie:latest`. Then run with `docker run -it collie:latest --CleanBranch="feature/UNI-1234-test"`

### Private Docker
TODO:

## Running Tests
Test commands should be ran from the `lib` directory. Go has several commands for testing. Test files in Go are prepended with `_test.go`. Inside test files a test func must begin with `Test`. Go also has Benchmark tests built in. A benchmark func must begin with `Benchmark`. Benchmarks are useful to see how a change affects performance.

- `go test`: Runs all tests.
- `go test -cover`: Runs all tests and gives coverage.
- `go test -bench .`: Runs all benchmarks.

### Versioning

## Contributing

## Related Projects

## Resources