## corectl context login

Login and set cookie for the named context

### Synopsis

Login and set cookie for the named context
	
This is only applicable when connecting to 'Qlik Sense Enterprise for Windows' through its proxy using HTTPS.
If no 'context-name' is used as argument the 'current-context' defined in the config will be used instead.

```
corectl context login <context-name> [flags]
```

### Examples

```
corectl context login
corectl context login context-name
```

### Options

```
  -h, --help              help for login
      --password string   Password to be used when logging in to Qlik Sense Enterprise (use with caution)
      --user string       Username to be used when logging in to Qlik Sense Enterprise
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Name of the context used when connecting to Qlik Associative Engine
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl context](corectl_context.md)	 - Create, update and use contexts

