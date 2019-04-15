[![CircleCI](https://circleci.com/gh/qlik-oss/corectl.svg?style=shield)](https://circleci.com/gh/qlik-oss/corectl)
[![Go Report Card](https://goreportcard.com/badge/qlik-oss/corectl)](https://goreportcard.com/report/qlik-oss/corectl)
![Latest Version](https://img.shields.io/github/release/qlik-oss/corectl.svg?style=flat)

<img src="./corectl.svg" alt="corectl" width="200"/>

## (Experimental)
Corectl is a command line tool to perform reloads, fetch metadata and evaluate expressions in Qlik Core apps.

---

## Download

On **Linux** and **OS X**

```bash
 curl --silent --location "https://github.com/qlik-oss/corectl/releases/latest/download/corectl-$(uname -s)-x86_64.tar.gz" | tar xz -C /tmp && mv /tmp/corectl /usr/local/bin/corectl
```

On **Windows** with git bash

```bash
curl --silent --location "https://github.com/qlik-oss/corectl/releases/latest/download/corectl-windows-x86_64.zip" > corectl.zip && unzip ./corectl.zip -d "$HOME/bin/" && rm ./corectl.zip
```

You can also download the binary manually from [releases](https://github.com/qlik-oss/corectl/releases).

## Examples

This sections describes some commands and configuration that can be used with the `corectl` tool.

To simplify usage of `corectl`, basic configurations such as: engine connection details, app and objects, can be described in a configuration file.
We have added an example configuration file to this repo [here](./examples/corectl.yml).

`corectl` will automatically check for a `corectl.yml | corectl.yaml` file in your current directory, removing the need to pass the config file using flags for each command.

Example configuration:
```yaml
engine: localhost:9076 # URL and port to running Qlik Associative Engine instance
app: corectl-example.qvf # App name that the tool should open a session against. Default a session app will be used.
script: ./script.qvs # Path to a script that should be set in the app
connections: # Connections that should be created in the app
  testdata: # Name of the connection
    connectionstring: /data # Connectionstring (qConnectionString) of the connection.
    type: folder # Type of connection
```

For more information regarding which additional options that are configurable are further described [here](./docs/corectl_config.md).

![](./examples/corectl-example.gif)

Also check out the blog post about utilizing `corectl` and `catwalk` to build your data model [here](https://branch-blog.qlik.com/data-modelling-in-qlik-core-a2e657c7598d).

## Usage

Usage documentation can be found [here](./docs/corectl.md).

`corectl` provides auto completion of commands and flags for `bash` and `zsh`. To load completion in your shell add the following to your `~/.bashrc` or `~/.zshrc` file depending on shell.

`. <(corectl completion bash)` or `. <(corectl completion zsh)`

Auto completion requires `bash-completion` to be installed.


# Development

## Prerequisite
- golang >= 1.11

## Build

Fast and easy - corectl will be built into the `$GOPATH/bin` and executable directly from bash using `corectl`
```bash
go install
```

If you want to keep the previous installed version you can use `go build` and get the binary to the current working directory
```bash
go build
```

## Test

The unit tests are run with the go test command:

```sh
go test ./...
```

The integration tests depend on external components. Before they can run, you must accept the [Qlik Core EULA](https://core.qlik.com/eula/) 
by setting the `ACCEPT_EULA` environment variable, you start the services by using the [docker-compose.yml](./test/docker-compose.yml) file.
The tests are run with the test script:

```sh
ACCEPT_EULA=<yes/no> docker-compose up -d
go test corectl_integration_test.go
```

The tests are by default trying to connect to an engine on localhost:9076 and another one on localhost:9176. This can be changed with the --engineIP flag and --engine2IP flag.

```sh
go test corectl_integration_test.go --engineIP HOST:PORT --engine2IP HOST:ANOTHERPORT
```

If the reference output files need to be updated, run the test with --update flag.

```sh
go test corectl_integration_test.go --update
```

## Release

You create a release by pushing a git tag with semantic versioning.
CircleCi will then run a release build that uses `goreleaser` to release `corectl` with the version set as the git tag.

Example:
`git tag v0.1.0 COMMIT-SHA`
`git push origin v0.1.0`

## Documentation

The usage documentation is generated using [`cobra/doc`](https://github.com/spf13/cobra/blob/master/doc/md_docs.md).
To regenerate the documentation:

```bash
corectl generate-docs
```

To regenerate the api spec, first build with latest release
tag as version and then generate the spec using:
```bash
go build -ldflags "-X main.version=$(git describe --abbrev=0 --tags)"
./corectl generate-spec > docs/spec.json
```

## Contributing

We welcome and encourage contributions! Please read [Open Source at Qlik R&D](https://github.com/qlik-oss/open-source)
for more info on how to get involved.
