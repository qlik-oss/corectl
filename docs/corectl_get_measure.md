## corectl get measure

Shows content of an generic measure

### Synopsis

Shows content of an generic measure. If no subcommand is specified the properties will be shown. Example: corectl get measure MEASURE-ID --app my-app.qvf

```
corectl get measure <measure-id> [flags]
```

### Options

```
  -a, --app string   App name, if no app is specified a session app is used instead.
  -h, --help         help for measure
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -j, --json                     Set logging format to JSON
  -l, --log-level string         Set logging level, one of; TRACE, DEBUG, INFO, WARN, ERROR, FATAL and PANIC. Logging levels DEBUG and TRACE includes JSON websocket traffic. (default "INFO")
      --ttl string               Engine session time to live in seconds (default "30")
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources
* [corectl get measure layout](corectl_get_measure_layout.md)	 - Evaluates the layout of an generic measure
* [corectl get measure properties](corectl_get_measure_properties.md)	 - Prints the properties of the generic measure

