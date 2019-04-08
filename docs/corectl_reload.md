## corectl reload

Reloads the app.

### Synopsis

Reloads the app.

```
corectl reload [flags]
```

### Examples

```
corectl reload
```

### Options

```
      --connections string   Path to a yml file containing the data connection definitions
      --dimensions string    A list of generic dimension json paths
  -h, --help                 help for reload
      --measures string      A list of generic measures json paths
      --no-save              Do not save the app
      --objects string       A list of generic object json paths
      --script string        path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded
      --silent               Do not log reload output
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

* [corectl](corectl.md)	 - 

