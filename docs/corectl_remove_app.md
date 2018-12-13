## corectl remove app

removes the specified app.

### Synopsis

removes the specified app.

```
corectl remove app <app-id> [flags]
```

### Examples

```
corectl remove app APP-ID
```

### Options

```
  -h, --help   help for app
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl remove](corectl_remove.md)	 - Remove entities (connections, dimensions, measures, objects) in the app or the app itself

