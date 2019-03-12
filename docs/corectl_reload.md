## corectl reload

Reloads the app.

### Synopsis

Reloads the app. Example: corectl reload

```
corectl reload [flags]
```

### Options

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to engine (default "localhost:9076")
      --headers stringToString   Headers to use when connecting to qix engine (default [])
  -h, --help                     help for reload
      --no-save                  Do not save the app
      --silent                   Do not log reload progress
      --ttl string               Engine session time to live in seconds (default "30")
```

### Options inherited from parent commands

```
  -j, --json               Set logging format to JSON
  -l, --log-level string   Set logging level, one of; TRACE, DEBUG, INFO, WARN, ERROR, FATAL and PANIC. Logging levels DEBUG and TRACE includes JSON websocket traffic. (default "INFO")
```

### SEE ALSO

* [corectl](corectl.md)	 - 

