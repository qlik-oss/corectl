## corectl get object data

Evaluates the hypercube data of an generic object

### Synopsis

Evaluates the hypercube data of an generic object. Example: corectl get object data OBJECT-ID --app my-app.qvf

```
corectl get object data <object-id> [flags]
```

### Options

```
  -h, --help   help for data
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

* [corectl get object](corectl_get_object.md)	 - Shows content of an generic object

