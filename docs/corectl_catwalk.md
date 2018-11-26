## corectl catwalk

Opens the specified app in Catwalk

### Synopsis

Opens the app in Catwalk. Example: corectl catwalk --app my-app.qvf
			


```
corectl catwalk [flags]
```

### Options

```
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for catwalk
      --ttl string      Engine session time to live in seconds (default "30")
```

### Options inherited from parent commands

```
  -v, --verbose   Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

