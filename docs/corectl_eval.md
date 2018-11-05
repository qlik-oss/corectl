## corectl eval

Evalutes a list of measures and dimensions

### Synopsis

Evalutes a list of measures and dimensions. To evaluate a measure for a specific dimension use the <measure> by <dimension> notation. If dimensions are omitted then the eval will be evaluated over all dimensions.

```
corectl eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...] [flags]
```

### Examples

```
corectl eval Count(a) // returns the number of values in field "a"
corectl eval 1+1 // returns the calculated value for 1+1
corectl eval Avg(Sales) by Region // returns the average of measure "Sales" for dimension "Region"
```

### Options

```
  -h, --help   help for eval
```

### Options inherited from parent commands

```
  -a, --app string               App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

