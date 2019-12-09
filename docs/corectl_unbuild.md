## corectl unbuild

Split up an existing app into separate json and yaml files

### Synopsis

Extracts generic objects, dimensions, measures, variables, reload script and connections from an app in an engine into separate json and yaml files.
In addition to the resources from the app a corectl.yml configuration file is generated that binds them all together.
Passwords in the connection definitions can not be exported from the app and hence need to be handled manually.
Generic Object trees (e.g. Qlik Sense sheets) are exported as a full property tree which means that child objects are found inside the parent´s json (the qChildren array).


```
corectl unbuild [flags]
```

### Examples

```
corectl unbuild
corectl unbuild --app APP-ID
```

### Options

```
      --dir string   Path to a the folder where the unbuilt app is exported (default "./<app name>-unbuild")
  -h, --help         help for unbuild
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Name of the context used when connecting to Qlik Associative Engine
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

