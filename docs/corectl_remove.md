## corectl remove

remove one or mores generic entities (dimensions, measures, objects) in the app

### Synopsis

remove one or mores generic entities (dimensions, measures, objects) in the app

### Options

```
  -a, --app string               App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -h, --help                     help for remove
      --noSave                   Do not save the app after doing reload
      --ttl string               Engine session time to live (default "30")
```

### Options inherited from parent commands

```
  -v, --verbose   Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl remove dimensions](corectl_remove_dimensions.md)	 - removes the specified generic dimensions in the current app
* [corectl remove measures](corectl_remove_measures.md)	 - removes the specified generic measures in the current app
* [corectl remove objects](corectl_remove_objects.md)	 - removes the specified generic objects in the current app

