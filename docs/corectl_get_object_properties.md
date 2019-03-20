## corectl get object properties

Prints the properties of the generic object

### Synopsis

Prints the properties of the generic object in JSON format

```
corectl get object properties <object-id> [flags]
```

### Examples

```
corectl get object properties OBJECT-ID --app my-app.qvf
```

### Options

```
  -h, --help   help for properties
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

* [corectl get object](corectl_get_object.md)	 - Shows content of an generic object

