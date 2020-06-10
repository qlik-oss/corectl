## corectl context

Create, update and use contexts

### Synopsis

Create, update and use contexts

Contexts store connection information such as server url, certificates and headers,
similar to a config. The main difference between contexts and configs is that they
can be used globally. Use the context subcommands to configure contexts which
facilitate app development in environments where certificates and headers are needed.

The current context is the one that is being used. You can use "context get" to
display the contents of the current context and switch context with "context set"
or unset the current context with "context unset".

Note that contexts have the lowest precedence. This means that e.g. an --server flag
(or a server field in a config) will override the server url in the current context.

Contexts are stored locally in your ~/.qlik/contexts.yml file.

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
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -s, --server string            URL to a Qlik Product, a local engine, cluster or sense-enterprise
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl context clear](corectl_context_clear.md)	 - Set the current context to none
* [corectl context create](corectl_context_create.md)	 - Create a context with the specified configuration
* [corectl context get](corectl_context_get.md)	 - Get context, current context by default
* [corectl context init](corectl_context_init.md)	 - Set up access to Qlik Sense SaaS
* [corectl context login](corectl_context_login.md)	 - Login and set cookie for the named context
* [corectl context ls](corectl_context_ls.md)	 - List all contexts
* [corectl context rm](corectl_context_rm.md)	 - Remove one or more contexts
* [corectl context update](corectl_context_update.md)	 - Update a context with the specified configuration
* [corectl context use](corectl_context_use.md)	 - Specify what context to use

