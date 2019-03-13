## corectl get dimensions

Prints a list of all generic dimensions in the current app

### Synopsis

Prints a list of all generic dimensions in the current app

```
corectl get dimensions [flags]
```

### Options

```
  -a, --app string   App name, if no app is specified a session app is used instead.
  -h, --help         help for dimensions
      --json         Prints the information in json format
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources

