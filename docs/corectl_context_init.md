---
title: "corectl context init"
description: "corectl context init"
categories: Libraries & Tools
type: Tools
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl context init

Set up access to Qlik Sense Cloud

### Synopsis

Set up access to Qlik Sense on Cloud Services/Kubernetes by entering the domain name and the api key of the Qlik Sense instance. If no context name is supplied the domain name is used as context name

```
corectl context init <context name> [flags]
```

### Options

```
      --api-key string   API key of the tenant
  -h, --help             help for init
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

* [corectl context](/commands/corectl_context)	 - Create, update and use contexts

