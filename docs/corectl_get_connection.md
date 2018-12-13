## corectl get connection

Shows the properties for a specific connection

### Synopsis

Shows the properties for a specific connection

```
corectl get connection [flags]
```

### Examples

```
corectl get connection CONNECTION-ID
```

### Options

```
  -a, --app string   App name, if no app is specified a session app is used instead.
  -h, --help         help for connection
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources

