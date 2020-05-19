---
title: "corectl dimension"
description: "corectl dimension"
categories: Libraries & Tools
type: Commands
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl dimension

Explore and manage dimensions

### Synopsis

Explore and manage dimensions

### Options

```
  -h, --help   help for dimension
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

* [corectl](/libraries-and-tools/corectl)	 - 
* [corectl dimension layout](/libraries-and-tools/corectl-dimension-layout)	 - Evaluate the layout of an generic dimension
* [corectl dimension ls](/libraries-and-tools/corectl-dimension-ls)	 - Print a list of all generic dimensions in the current app
* [corectl dimension properties](/libraries-and-tools/corectl-dimension-properties)	 - Print the properties of the generic dimension
* [corectl dimension rm](/libraries-and-tools/corectl-dimension-rm)	 - Remove one or many dimensions in the current app
* [corectl dimension set](/libraries-and-tools/corectl-dimension-set)	 - Set or update the dimensions in the current app

