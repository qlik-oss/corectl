## corectl script script

Print the reload script

### Synopsis

Fetches the script currently set in the app and prints it in plain text.

```
corectl script script [flags]
```

### Examples

```
corectl get script
corectl get script --app=my-app.qvf
```

### Options

```
  -h, --help   help for script
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

* [corectl script](corectl_script.md)	 - Explore and manage the script

