[![CircleCI](https://circleci.com/gh/qlik-oss/corectl.svg?style=shield)](https://circleci.com/gh/qlik-oss/corectl)
[![Go Report Card](https://goreportcard.com/badge/qlik-oss/corectl)](https://goreportcard.com/report/qlik-oss/corectl)
![Latest Version](https://img.shields.io/github/release/qlik-oss/corectl.svg?style=flat)

# Corectl (Experimental)

Corectl is a command line tool to perform reloads, fetch metadata and evaluate expressions in Qlik Core apps.

---

## Download

**Change _\<version\>_ below to the version you want to download.**

E.G v0.0.4

On **Linux** and **OS X**

```bash
 curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/<version>/corectl-$(uname -s)-x86_64.tar.gz" | tar xz -C /tmp && mv /tmp/corectl /usr/local/bin/corectl
```

On **Windows** with git bash

```bash
curl --silent --location "https://github.com/qlik-oss/corectl/releases/download/<version>/corectl-windows-x86_64.zip" > corectl.zip && unzip ./corectl.zip -d "$HOME/bin/" && rm ./corectl.zip
```

You can also download the binary manually from [releases](https://github.com/qlik-oss/corectl/releases).

## Development

Either clone the repo or go get it:

```bash
go get -u github.com/qlik-oss/corectl
```

Build the main.go file to a location on your path. You can use the buildtohomebin script.

```bash
./buildtohomebin
```

## Usage

Usage documentation and examples can be found [here](./docs/corectl.md).

`corectl` provides autocompletion of commands for `bash` and `zsh`. To load completion in your shell add the following to your `~/.bashrc` or `~/.zshrc` file depending on shell.

`. <(corectl completion bash)` or `. <(corectl completion zsh)`

Autocompletion requires `jq` and `bash-completion` to be installed.

## Testing

The unit tests are run with the go test command:

```sh
go test ./...
```

The integration tests depend on external components. Before they can run, you must accept the [Qlik Core EULA](https://core.qlik.com/eula/) 
by setting the `ACCEPT_EULA` environment variable, you start the services by using the [docker-compose.yml](./docker-compose.yml) file.
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

`corectl generate-docs`

## Contributing

We welcome and encourage contributions! Please read [Open Source at Qlik R&D](https://github.com/qlik-oss/open-source)
for more info on how to get involved.
