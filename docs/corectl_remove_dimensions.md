## corectl remove dimensions

removes the specified generic dimensions in the current app

### Synopsis

removes the specified generic dimensions in the current app. Example: corectl remove dimension ID-1 ID-2

```
corectl remove dimensions [flags]
```

### Options

```
  -h, --help   help for dimensions
```

### Options inherited from parent commands

```
  -a, --app string               App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --noSave                   Do not save the app after doing reload
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl remove](corectl_remove.md)	 - remove one or mores generic entities (dimensions, measures, objects) in the app

