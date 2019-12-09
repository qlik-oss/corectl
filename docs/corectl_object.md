## corectl object

Explore and manage generic objects

### Synopsis

Explore and manage generic objects

### Options

```
  -h, --help   help for object
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
* [corectl object data](corectl_object_data.md)	 - Evaluate the hypercube data of a generic object
* [corectl object layout](corectl_object_layout.md)	 - Evaluate the hypercube layout of the generic object
* [corectl object ls](corectl_object_ls.md)	 - Print a list of all generic objects in the current app
* [corectl object properties](corectl_object_properties.md)	 - Print the properties of the generic object
* [corectl object rm](corectl_object_rm.md)	 - Remove one or many generic objects in the current app
* [corectl object set](corectl_object_set.md)	 - Set or update the objects in the current app

