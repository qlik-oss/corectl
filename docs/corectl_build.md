## corectl build

Reload and save the app after updating connections, dimensions, measures, objects and the script

### Synopsis

Reload and save the app after updating connections, dimensions, measures, objects and the script

```
corectl build [flags]
```

### Examples

```
corectl build
corectl build --connections ./myconnections.yml --script ./myscript.qvs
```

### Options

```
      --bookmarks string     A list of generic bookmark json paths
      --connections string   Path to a yml file containing the data connection definitions
      --dimensions string    A list of generic dimension json paths
  -h, --help                 help for build
      --measures string      A list of generic measures json paths
      --no-reload            Do not run the reload script
      --no-save              Do not save the app
      --objects string       A list of generic object json paths
      --script string        Path to a qvs file containing the app data reload script
      --silent               Do not log reload output
      --variables string     A list of generic variable json paths
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -q, --quiet                    Terse logging, limited to printing IDs where possible
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

