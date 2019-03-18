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

## Examples

This sections describes some normal uses cases and configuration that can be used with the `corectl` tool.

To simply the usage of `corectl` the basic details for connecting to an engine instance and which apps or objects to interact with can be configured in a configuration file.
We have added an example configuration file to this repo [here](./examples/corect.yml).

`corectl` will automatically check for a `corectl.yml | corectl.yaml` file in your current directory. It is also possible to specify a configuration file using the `--config` or `-c` flag.

Example configuration:
```yaml
engine: localhost:9076 # URL and port to running Qlik Associative Engine instance
app: corectl-example.qvf # App name that the tool should open a session against. Default a session app will be used.
script: ./script.qvs # Path to a script that should be set in the app
connections: # Connections that should be created in the app
  testdata: # Name of the connection
    connectionstring: /data # Connectionstring (qConnectionString) of the connection. For a folder connector this is an absolute or relative path inside of the engine docker container.
    type: folder # Type of connection
objects:
  - ./object-*.json # Path to objects that should be created from a json file. Accepts wildcards.
```

All of the configurations that are possible to specify in a configuration file is also possible to pass by command line flags. You can find all the supported commands and flags in the [Usage section](#usage).

![](./examples/corectl_example.gif)

Also check out the blog post about utilizing `corectl` and `catwalk` to build your data model [here](https://branch-blog.qlik.com/data-modelling-in-qlik-core-a2e657c7598d).

## Usage

Usage documentation can be found [here](./docs/corectl.md).

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
