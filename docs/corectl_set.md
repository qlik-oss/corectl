## corectl set

Sets one or several resources

### Synopsis

Sets one or several resources

### Options

```
  -h, --help      help for set
      --no-save   Do not save the app
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl set all](corectl_set_all.md)	 - Sets the objects, measures, dimensions, connections and script in the current app
* [corectl set connections](corectl_set_connections.md)	 - Sets or updates the connections in the current app
* [corectl set dimensions](corectl_set_dimensions.md)	 - Sets or updates the dimensions in the current app
* [corectl set measures](corectl_set_measures.md)	 - Sets or updates the measures in the current app
* [corectl set objects](corectl_set_objects.md)	 - Sets or updates the objects in the current app
* [corectl set script](corectl_set_script.md)	 - Sets the script in the current app

