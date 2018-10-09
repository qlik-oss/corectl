## corectl eval

Evalutes a list of measures and dimensions

### Synopsis

Evalutes a list of measures and dimensions. Meaures are separeted from dimensions by the "by" keyword. To omit dimensions and only use measures use "*" as dimension: eval <measures> by *

```
corectl eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...] [flags]
```

### Options

```
  -h, --help   help for eval
```

### Options inherited from parent commands

```
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
      --ttl string      Engine session time to live (default "30")
  -v, --verbose         Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

