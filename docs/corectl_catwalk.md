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
  -a, --app string           App name, if no app is specified a session app is used instead.
      --catwalk-url string   Url to an instance of catwalk, if not provided the qlik one will be used. (default "https://catwalk.core.qlik.com")
  -c, --config string        path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string        URL to engine (default "localhost:9076")
  -h, --help                 help for catwalk
```

### Options inherited from parent commands

```
  -j, --json               Set logging format to JSON
  -l, --log-level string   Set logging level, one of; TRACE, DEBUG, INFO, WARN, ERROR, FATAL and PANIC. Logging levels DEBUG and TRACE includes JSON websocket traffic. (default "INFO")
```

### SEE ALSO

* [corectl](corectl.md)	 - 

