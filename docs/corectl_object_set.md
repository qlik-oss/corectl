## corectl object set

Sets or updates the objects in the current app

### Synopsis

Sets or updates the objects in the current app.
The JSON objects can be in either the GenericObjectProperties format or the GenericObjectEntry format

```
corectl object set <glob-pattern-path-to-objects-files.json [flags]
```

### Examples

```
corectl object set ./my-objects-glob-path.json
```

### Options

```
  -h, --help   help for set
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-data                  Open app without data
      --no-save                  Do not save the app
      --suppress                 Suppress all confirmation dialogues
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl object](corectl_object.md)	 - Explore and manage generic objects

