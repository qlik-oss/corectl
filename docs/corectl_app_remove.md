## corectl app remove

Remove the specified app.

### Synopsis

Remove the specified app.

```
corectl app remove <app-id> [flags]
```

### Examples

```
corectl app remove
corectl app remove APP-ID
```

### Options

```
  -h, --help       help for remove
      --suppress   Suppress confirmation dialogue
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl app](corectl_app.md)	 - Explore and manage apps

