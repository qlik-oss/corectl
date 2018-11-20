## corectl set objects

Sets or updates the objects in the current app

### Synopsis

Sets or updates the objects in the current app Example corectl set objects ./my-objects-glob-path.json

```
corectl set objects [flags]
```

### Options

```
  -h, --help   help for objects
```

### Options inherited from parent commands

```
  -a, --app string               App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --noSave                   Do not save the app
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl set](corectl_set.md)	 - Sets one or several resources

