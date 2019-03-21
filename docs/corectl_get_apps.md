## corectl get apps

Prints a list of all apps available in the current engine

### Synopsis

Prints a list of all apps available in the current engine

```
corectl get apps [flags]
```

### Examples

```
corectl get apps
corectl get apps --engine=localhost:9276
```

### Options

```
  -h, --help   help for apps
      --json   Prints the information in json format
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl get](corectl_get.md)	 - Lists one or several resources

