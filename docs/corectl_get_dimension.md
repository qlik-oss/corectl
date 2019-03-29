## corectl get dimension

Shows content of an generic dimension

### Synopsis

Shows content of an generic dimension. If no subcommand is specified the properties will be shown.

```
corectl get dimension <dimension-id> [flags]
```

### Examples

```
corectl get dimension DIMENSION-ID --app my-app.qvf
```

### Options

```
  -h, --help   help for dimension
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
* [corectl get dimension layout](corectl_get_dimension_layout.md)	 - Evaluates the layout of an generic dimension
* [corectl get dimension properties](corectl_get_dimension_properties.md)	 - Prints the properties of the generic dimension

