## corectl reload

Reloads the app

### Synopsis

Reloads the app. Example: corectl reload --connections ./myconnections.yml --script ./myscript.qvs

```
corectl reload [flags]
```

### Options

```
      --connections string   path to connections file
  -h, --help                 help for reload
      --script string        Script file name
```

### Options inherited from parent commands

```
  -a, --app string              App name including .qvf file ending
  -c, --config string           path/to/config.yml where default parameters can be set
  -e, --engine string           URL to engine
      --engine-headers string   HTTP headers to send to the engine (default "30")
      --ttl string              Engine session time to live (default "30")
  -v, --verbose                 Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

