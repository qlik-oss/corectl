## corectl remove connection

Remove the specified connection(s)

### Synopsis

Remove one or many connections from the app

```
corectl remove connection <connection-id>... [flags]
```

### Examples

```
corectl remove connection ID-1
corectl remove connections ID-1 ID-2
```

### Options

```
  -h, --help      help for connection
      --no-save   Do not save the app
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-data                  Open app without data
      --suppress                 Suppress all confirmation dialogues
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl remove](corectl_remove.md)	 - Remove entities (connections, dimensions, measures, objects) in the app or the app itself

