## corectl catwalk

Opens the specified app in catwalk

### Synopsis

Opens the specified app in catwalk. If no app is specified the catwalk hub will be opened.

```
corectl catwalk [flags]
```

### Examples

```
corectl catwalk --app my-app.qvf
corectl catwalk --app my-app.qvf --catwalk-url http://localhost:8080
```

### Options

```
      --catwalk-url string   Url to an instance of catwalk, if not provided the qlik one will be used. (default "https://catwalk.core.qlik.com")
  -h, --help                 help for catwalk
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

* [corectl](corectl.md)	 - 

