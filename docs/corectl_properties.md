## corectl properties

Prints the properties of the object identified by the --object flag and the type with --type

### Synopsis

Prints the properties of the object identified by the --object flag and the type with --type. If --type is ommited genericObject is assumed

```
corectl properties [flags]
```

### Options

```
  -h, --help              help for properties
      --measures string   A list of measure json paths
  -o, --object string     ID of a generic object
      --objects string    A list of object json paths
      --type string       The type of object to print
```

### Options inherited from parent commands

```
  -a, --app string               App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

