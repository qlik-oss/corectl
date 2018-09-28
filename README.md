# Corectl (Experimental)
Corectl is a command line tool to perform reloads, fetch metadata and evaluate expressions in Qlik Core apps.


---

## Download

On OS X

```bash
curl -L https://github.com/qlik-oss/corectl/releases/download/v0.0.3/corectl-darwin-amd64 >/usr/local/bin/corectl &&
  chmod +x /usr/local/bin/corectl
```

On Linux

```bash
curl -L https://github.com/qlik-oss/corectl/releases/download/v0.0.3/corectl-linux-amd64 >/tmp/corectl &&
    chmod +x /tmp/corectl &&
   cp /tmp/corectl /usr/local/bin/corectl
```

On Windows with git bash

```bash
mkdir -p "$HOME/bin" &&
curl -L https://github.com/qlik-oss/corectl/releases/download/v0.0.3/corectl-windows-amd64.exe > "$HOME/bin/corectl.exe" &&
chmod +x "$HOME/bin/corectl.exe"
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

## Testing

The unit tests are run with the go test command:

```sh
$ go test ./...
```

The integration tests depend on external components. Before they can run, you must accept the [Qlik Core EULA](https://core.qlik.com/eula/) 
by setting the `ACCEPT_EULA` environment variable, you start the services by using the [docker-compose.yml](./docker-compose.yml) file.
The tests are run with the test script:


```sh
$ ACCEPT_EULA=<yes/no> docker-compose up -d
$ go test corectl_integration_test.go
```

The tests are by default trying to connect to an engine on localhost:9076. This can be changed with the --engineIP flag.

```sh
$ go test corectl_integration_test.go --engineIP HOST:PORT
```

If the reference output files need to be updated, run the test with --update flag.

```sh
$ go test corectl_integration_test.go --update
```

## Documentation

The usage documentation is generated using [`cobra/doc`](https://github.com/spf13/cobra/blob/master/doc/md_docs.md).
To regenerate the documentation:

`corectl generate-docs`

## Contributing
We welcome and encourage contributions! Please read [Open Source at Qlik R&D](https://github.com/qlik-oss/open-source)
for more info on how to get involved.
