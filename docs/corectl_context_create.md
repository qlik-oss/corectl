## corectl context create

Create a new context

### Synopsis

Create a new context

```
corectl context create <context name> [flags]
```

### Examples

```
corectl context create local-engine
corectl context create rd-sense --product "QSE" --comment "R&D Qlik Sense deployment"
```

### Options

```
      --comment string   Comment for the context
  -h, --help             help for create
      --product string   Qlik product the context is connecting to. One of: QC (Qlik Core), QSC (Qlik Sense Cloud), QSD (Qlik Sense Desktop), QSE (Qlik Sense Enterprise), QSEoK (Qlik Sense Enterpries on Kubernetes) or QSEoW (Qlik Sense Enterpries on Windows) (default "QC")
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Specific context that should be used when connecting
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl context](corectl_context.md)	 - Explore and manage contexts

