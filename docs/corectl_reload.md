## corectl reload

Reloads and saves the app after updating connections, objects and the script

### Synopsis

Reloads the app. Example: corectl reload --connections ./myconnections.yml --script ./myscript.qvs
			


```
corectl reload [flags]
```

### Examples

```
  # Specify all parameters on the command line:
  corectl reload --connections ./myconnections.yml --script ./myscript.qvs

  # Specify parameters in the config file:
  corectl reload --config ./config.yml

```

### Options

```
      --connections string   path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they are specified inline, not as a file reference!
  -h, --help                 help for reload
      --objects string       A list of object json paths
      --script string        path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded
      --silent               Do not log reload progress
```

### Options inherited from parent commands

```
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
      --ttl string      Engine session time to live (default "30")
  -v, --verbose         Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

