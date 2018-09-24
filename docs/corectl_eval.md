## corectl eval

Evalutes a hypercube

### Synopsis

Evalutes a list of measures and dimensions. Meaures are separeted from dimensions by the "by" keyword. To omit dimensions and only use measures use "*" as dimension: eval <measures> by *

```
corectl eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...] [flags]
```

### Options

```
  -h, --help            help for eval
  -s, --select string   
```

### Options inherited from parent commands

```
  -a, --app string              App name including .qvf file ending
  -c, --config string           path/to/config.yml where default parameters can be set
  -e, --engine string           URL to engine
      --engine-headers string   HTTP headers to send to the engine (default "30")
      --ttl string              Engine session time to live (default "30")
  -v, --verbose                 Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

