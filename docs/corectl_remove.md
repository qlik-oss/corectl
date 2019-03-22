## corectl remove

Remove entities (connections, dimensions, measures, objects) in the app or the app itself

### Synopsis

Remove one or mores generic entities (connections, dimensions, measures, objects) in the app

### Options

```
  -h, --help       help for remove
      --suppress   Suppress all confirmation dialogues
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl remove app](corectl_remove_app.md)	 - removes the specified app.
* [corectl remove connections](corectl_remove_connections.md)	 - Removes the specified connection(s)
* [corectl remove dimensions](corectl_remove_dimensions.md)	 - Removes the specified generic dimensions in the current app
* [corectl remove measures](corectl_remove_measures.md)	 - Removes the specified generic measures in the current app
* [corectl remove objects](corectl_remove_objects.md)	 - Removes the specified generic objects in the current app

