## corectl get field

Shows content of a field

### Synopsis

Prints all the values for a specific field in your data model

```
corectl get field <field name> [flags]
```

### Examples

```
corectl get field FIELD
```

### Options

```
  -a, --app string   App name, if no app is specified a session app is used instead.
  -h, --help         help for field
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources

