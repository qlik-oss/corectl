# Corectl (Experimental)
Corectl is a command line tool to perform reloads, fetch metadata and evaluate expressions in Qlik Core apps.


---

## Installation
Either clone the repo or go get it:
```bash
go get -u github.com/qlik-oss/corectl
```

Build the main.go file to a location on your path. You can use the buildtohomebin script.
```bash
./builtohomebin
```

## Example Usage
Reload a script file in the specified app and print metadata. The script file path is local, the app name/path is from within the engine docker file system.
```bash
corectl --app myapp.qvf reload myscript.qvs
```

Print the metadata with reload
```bash
corectl --app myapp.qvf meta
```

Evaluate expressions. Note the "by" keyword. The format is `<expressions> by <dimensions>`.

```bash
corectl --app myapp.qvf eval "sum(Z)" by X Y
```

or iterate over all dimensions:

```bash
corectl --app myapp.qvf eval "sum(Z)" by "*"
```

The `eval` command can also be used for calculated dimensions:

```bash
corectl --app myapp.qvf eval "=A+B+C"
```

Specify what Qlik Associative Engine to use with the --engine parameter
```bash
corectl --engine remoteengine:9076 --app myapp.qvf reload myscript.qvs
```

Print some extra debugging information using --verbose flag
```bash
corectl --verbose --app myapp.qvf meta
```

## Testing

The unit tests are run with the go test command:

```sh
$ go test ./...
```

The integration tests depend on external components. Before they can run, you must accept the [Qlik Core EULA](https://qlikcore.com/beta/) 
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

## Contributing
We welcome and encourage contributions! Please read [Open Source at Qlik R&D](https://github.com/qlik-oss/open-source)
for more info on how to get involved.
