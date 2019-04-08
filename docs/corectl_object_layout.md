## corectl object layout

Evaluates the hypercube layout of an generic object

### Synopsis

Evaluates the hypercube layout of an generic object in JSON format

```
corectl object layout <object-id> [flags]
```

### Examples

```
corectl object layout OBJECT-ID --app my-app.qvf
```

### Options

```
  -h, --help   help for layout
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

* [corectl object](corectl_object.md)	 - Explore and manage generic objects

