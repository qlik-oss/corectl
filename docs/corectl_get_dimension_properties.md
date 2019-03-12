## corectl get dimension properties

Prints the properties of the generic dimension

### Synopsis

Prints the properties of the generic dimension. Example: corectl get dimension properties DIMENSION-ID --app my-app.qvf

```
corectl get dimension properties <dimension-id> [flags]
```

### Options

```
  -h, --help   help for properties
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -j, --json                     Set logging format to JSON
  -l, --log-level string         Set logging level, one of; TRACE, DEBUG, INFO, WARN, ERROR, FATAL and PANIC. Logging levels DEBUG and TRACE includes JSON websocket traffic. (default "INFO")
      --ttl string               Engine session time to live in seconds (default "30")
```

### SEE ALSO

* [corectl get dimension](corectl_get_dimension.md)	 - Shows content of an generic dimension

