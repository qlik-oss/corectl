## corectl remove dimensions

Removes the specified generic dimensions in the current app

### Synopsis

Removes the specified generic dimensions in the current app. Example: corectl remove dimension ID-1 ID-2

```
corectl remove dimensions <dimension-id>... [flags]
```

### Options

```
  -h, --help      help for dimensions
      --no-save   Do not save the app
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --supress                  Supress confirm
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl remove](corectl_remove.md)	 - Remove entities (connections, dimensions, measures, objects) in the app or the app itself

