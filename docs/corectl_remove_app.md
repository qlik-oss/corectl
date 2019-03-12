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
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -j, --json                     Set logging format to JSON
  -l, --log-level string         Set logging level, one of; TRACE, DEBUG, INFO, WARN, ERROR, FATAL and PANIC. Logging levels DEBUG and TRACE includes JSON websocket traffic. (default "INFO")
      --suppress                 Suppress all confirmation dialogues
      --ttl string               Engine session time to live in seconds (default "30")
```

### SEE ALSO

* [corectl remove](corectl_remove.md)	 - Remove entities (connections, dimensions, measures, objects) in the app or the app itself

