## corectl dimension

Explore and manage dimensions

### Synopsis

Explore and manage dimensions

### Options

```
  -h, --help   help for dimension
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl dimension layout](corectl_dimension_layout.md)	 - Evaluate the layout of an generic dimension
* [corectl dimension ls](corectl_dimension_ls.md)	 - Print a list of all generic dimensions in the current app
* [corectl dimension properties](corectl_dimension_properties.md)	 - Print the properties of the generic dimension
* [corectl dimension rm](corectl_dimension_rm.md)	 - Remove one or many dimensions in the current app
* [corectl dimension set](corectl_dimension_set.md)	 - Set or update the dimensions in the current app
