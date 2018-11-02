## corectl objects

Prints a list of all objects in the current app

### Synopsis

Prints a list of all objects in the current app

```
corectl objects [flags]
```

### Options

```
  -h, --help              help for objects
      --measures string   A list of measure json paths
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

