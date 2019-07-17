## corectl variable

Explore and manage variables

### Synopsis

Explore and manage variables

### Options

```
  -h, --help   help for variable
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Specific context that should be used when connecting
  -e, --engine string            URL to the Qlik Associative Engine
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl variable layout](corectl_variable_layout.md)	 - Evaluate the layout of an generic variable
* [corectl variable ls](corectl_variable_ls.md)	 - Print a list of all generic variables in the current app
* [corectl variable properties](corectl_variable_properties.md)	 - Print the properties of the generic variable
* [corectl variable rm](corectl_variable_rm.md)	 - Remove one or many variables in the current app
* [corectl variable set](corectl_variable_set.md)	 - Set or update the variables in the current app

