## corectl connection

Explore and manage connections

### Synopsis

Explore and manage connections

### Options

```
  -h, --help   help for connection
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
* [corectl connection get](corectl_connection_get.md)	 - Show the properties for a specific connection
* [corectl connection ls](corectl_connection_ls.md)	 - Print a list of all connections in the current app
* [corectl connection rm](corectl_connection_rm.md)	 - Remove the specified connection(s)
* [corectl connection set](corectl_connection_set.md)	 - Set or update the connections in the current app

