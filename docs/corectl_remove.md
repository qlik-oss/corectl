## corectl remove

Remove entities (connections, dimensions, measures, objects) in the app or the app itself

### Synopsis

Remove one or mores generic entities (connections, dimensions, measures, objects) in the app

### Options

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -h, --help                     help for remove
      --suppress                 Suppress all confirmation dialogues
      --ttl string               Engine session time to live in seconds (default "30")
```

### Options inherited from parent commands

```
  -t, --traffic   Log JSON websocket traffic to stdout
  -v, --verbose   Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl remove app](corectl_remove_app.md)	 - removes the specified app.
* [corectl remove connection](corectl_remove_connection.md)	 - removes the specified connection.
* [corectl remove dimensions](corectl_remove_dimensions.md)	 - Removes the specified generic dimensions in the current app
* [corectl remove measures](corectl_remove_measures.md)	 - Removes the specified generic measures in the current app
* [corectl remove objects](corectl_remove_objects.md)	 - Removes the specified generic objects in the current app

