## corectl get object

Shows content of an generic object

### Synopsis

Shows content of an generic object. If no subcommand is specified the properties will be shown.

```
corectl get object <object-id> [flags]
```

### Examples

```
corectl get object OBJECT-ID --app my-app.qvf
```

### Options

```
  -h, --help   help for object
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

* [corectl get](corectl_get.md)	 - Lists one or several resources
* [corectl get object data](corectl_get_object_data.md)	 - Evaluates the hypercube data of an generic object
* [corectl get object layout](corectl_get_object_layout.md)	 - Evaluates the hypercube layout of an generic object
* [corectl get object properties](corectl_get_object_properties.md)	 - Prints the properties of the generic object

