## corectl get script

Print the reload script

### Synopsis

Print the reload script

```
corectl get script [flags]
```

### Options

```
  -a, --app string   App name including .qvf file ending. If no app is specified a session app is used instead.
  -h, --help         help for script
```

### Options inherited from parent commands

```
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources

