# Core CLI (Experimental)
Core CLI is a command line tool to perform reloads, fetch metadata and evaluate expressions in Qlik Core apps.


---

## Installation
Either clone the repo or go get it:
```bash
go get -u github.com/qlik-oss/core-cli
```

Build the main.go file to a location on your path. You can use the buildtohomebin script.
```bash
./builtohomebin
```

## Example Usage
Reload a script file in the specified app and print metadata. The script file path is local, the app name/path is from within the engine docker file system.
```bash
qli --app myapp.qvf reload myscript.qvs
```

Print the metadata with reload
```bash
qli --app myapp.qvf meta
```

Evaluate expressions. Note the "by" keyword. The format is <expressions> by <dimensions>.
```bash
qli --app myapp.qvf eval "sum(Z)" by X Y
```

Specify what Qlik Associative Engine to use with the --engine parameter
```bash
qli --engine remoteengine:9076 --app myapp.qvf reload myscript.qvs
```

Print some extra debugging information using --verbose flag
```bash
qli --verbose --app myapp.qvf meta
```

## Testing

The integration tests depend on external components. Before they can run, you must accept the [Qlik Core EULA](https://qlikcore.com/beta/) 
by setting the `ACCEPT_EULA` environment variable, you start the services by using the [docker-compose.yml](./docker-compose.yml) file.
The tests are run with the test script:


```sh
$ ACCEPT_EULA=<yes/no> docker-compose up -d
$ ./test.sh
```

## Contributing
We welcome and encourage contributions! Please read [Open Source at Qlik R&D](https://github.com/qlik-oss/open-source)
for more info on how to get involved.
