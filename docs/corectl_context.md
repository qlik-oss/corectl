## corectl context

Create, update and use contexts

### Synopsis

Create, update and use contexts

Contexts store connection information such as engine url, certificates and headers,
similar to a config. The main difference between contexts and configs is that they
can be used globally. Use the context subcommands to configure contexts which
facilitate app development in environments where certificates and headers are needed.

The current context is the one that is being used. You can use "context get" to
display the contents of the current context and switch context with "context set"
or unset the current context with "context unset".

Note that contexts have the lowest precedence. This means that e.g. an --engine flag
(or an engine field in a config) will override the engine url in the current context.

Contexts are stored locally in your ~/.corectl/contexts.yml file.

### Options

```
  -h, --help   help for context
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

* [corectl](corectl.md)	 - 
* [corectl context clear](corectl_context_clear.md)	 - Set the current context to none
* [corectl context get](corectl_context_get.md)	 - Get context, current context by default
* [corectl context login](corectl_context_login.md)	 - Login and set cookie for the named context
* [corectl context ls](corectl_context_ls.md)	 - List all contexts
* [corectl context rm](corectl_context_rm.md)	 - Remove one or more contexts
* [corectl context set](corectl_context_set.md)	 - Set a context to the current configuration
* [corectl context use](corectl_context_use.md)	 - Specify what context to use

