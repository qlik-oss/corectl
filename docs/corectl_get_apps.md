## corectl get apps

Prints a list of all apps available in the current engine

### Synopsis

Prints a list of all apps available in the current engine

```
corectl get apps [flags]
```

### Options

```
  -h, --help   help for apps
      --json   Prints the apps in json format
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

