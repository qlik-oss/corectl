## corectl get status

Prints status info about the connection to the engine and current app

### Synopsis

Prints status info about the connection to the engine and current app, and also the status of the data model

```
corectl get status [flags]
```

### Examples

```
corectl get status
corectl get status --app=my-app.qvf
```

### Options

```
  -h, --help   help for status
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources

