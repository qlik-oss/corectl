## corectl get measure properties

Prints the properties of the generic measure

### Synopsis

Prints the properties of the generic measure. Example: corectl get measure properties MEASURE-ID --app my-app.qvf

```
corectl get measure properties <measure-id> [flags]
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
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get measure](corectl_get_measure.md)	 - Shows content of an generic measure

