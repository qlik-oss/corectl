## corectl get dimension

Shows content of an generic dimension

### Synopsis

Shows content of an generic dimension. If no subcommand is specified the properties will be shown. Example: corectl get dimension DIMENSION-ID --app my-app.qvf

```
corectl get dimension [flags]
```

### Options

```
  -a, --app string   App name including .qvf file ending. If no app is specified a session app is used instead.
  -h, --help         help for dimension
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources
* [corectl get dimension layout](corectl_get_dimension_layout.md)	 - Evalutes the layout of an generic dimension
* [corectl get dimension properties](corectl_get_dimension_properties.md)	 - Prints the properties of the generic dimension

