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
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Name of the context used when connecting to Qlik Associative Engine
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl measure layout](corectl_measure_layout.md)	 - Evaluate the layout of an generic measure
* [corectl measure ls](corectl_measure_ls.md)	 - Print a list of all generic measures in the current app
* [corectl measure properties](corectl_measure_properties.md)	 - Print the properties of the generic measure
* [corectl measure rm](corectl_measure_rm.md)	 - Remove one or many generic measures in the current app
* [corectl measure set](corectl_measure_set.md)	 - Set or update the measures in the current app

