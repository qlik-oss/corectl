## corectl measure

Explore and manage measures

### Synopsis

Explore and manage measures

### Options

```
  -h, --help   help for measure
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
* [corectl measure layout](corectl_measure_layout.md)	 - Evaluate the layout of an generic measure
* [corectl measure ls](corectl_measure_ls.md)	 - Print a list of all generic measures in the current app
* [corectl measure properties](corectl_measure_properties.md)	 - Print the properties of the generic measure
* [corectl measure rm](corectl_measure_rm.md)	 - Remove one or many generic measures in the current app
* [corectl measure set](corectl_measure_set.md)	 - Set or update the measures in the current app

