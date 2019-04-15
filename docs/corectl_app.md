## corectl app

Explore and manage apps

### Synopsis

Explore and manage apps

### Options

```
  -h, --help   help for app
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
* [corectl app ls](corectl_app_ls.md)	 - Print a list of all apps available in the current engine
* [corectl app rm](corectl_app_rm.md)	 - Remove the specified app

