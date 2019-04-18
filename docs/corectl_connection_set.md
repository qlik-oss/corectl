## corectl connection set

Set or update the connections in the current app

### Synopsis

Set or update the connections in the current app

```
corectl connection set <path-to-connections-file.yml> [flags]
```

### Examples

```
corectl connection set ./my-connections.yml
```

### Options

```
  -h, --help   help for set
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --json                     Returns output in JSON format if possible, overrides the verbose and traffic flags
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl connection](corectl_connection.md)	 - Explore and manage connections

