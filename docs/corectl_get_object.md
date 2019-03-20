## corectl get object

Shows content of an generic object

### Synopsis

Shows content of an generic object. If no subcommand is specified the properties will be shown. Example: corectl get object OBJECT-ID --app my-app.qvf

```
corectl get object <object-id> [flags]
```

### Options

```
  -h, --help   help for object
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources
* [corectl get object data](corectl_get_object_data.md)	 - Evaluates the hypercube data of an generic object
* [corectl get object layout](corectl_get_object_layout.md)	 - Evaluates the hypercube layout of an generic object
* [corectl get object properties](corectl_get_object_properties.md)	 - Prints the properties of the generic object

