## corectl data

Evalutes the hypercube data of an object defined by the --object parameter. Note that only basic hypercubes like straight tables are supported

### Synopsis

Evalutes the hypercube data of an object defined by the --object parameter. Note that only basic hypercubes like straight tables are supported

```
corectl data [flags]
```

### Options

```
  -h, --help             help for data
  -o, --object string    ID of a generic object
      --objects string   A list of object json paths
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

