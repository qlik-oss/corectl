## corectl update

Updates connections, objects, measures and script and saves the app

### Synopsis

Updates connections, objects, measures and script in the app. Example: corectl update	



```
corectl update [flags]
```

### Options

```
      --connections string   path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they are specified inline, not as a file reference!
  -h, --help                 help for update
      --measures string      A list of measure json paths
      --objects string       A list of object json paths
      --script string        path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded
```

### Options inherited from parent commands

```
  -a, --app string               App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
      --ttl string               Engine session time to live (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

