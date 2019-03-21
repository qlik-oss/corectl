## corectl set connections

Sets or updates the connections in the current app

### Synopsis

Sets or updates the connections in the current app

```
corectl set connections <path-to-connections-file.yml> [flags]
```

### Examples

```
corectl set connections ./my-connections.yml
```

### Options

```
      --connections string   path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they are specified inline, not as a file reference!
  -h, --help                 help for connections
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-save                  Do not save the app
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl set](corectl_set.md)	 - Sets one or several resources

