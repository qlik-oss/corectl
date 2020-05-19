---
title: "corectl object"
description: "corectl object"
categories: Libraries & Tools
type: Commands
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl object

Explore and manage generic objects

### Synopsis

Explore and manage generic objects

### Options

```
  -h, --help   help for object
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
* [corectl object data](/libraries-and-tools/corectl-object-data)	 - Evaluate the hypercube data of a generic object
* [corectl object layout](/libraries-and-tools/corectl-object-layout)	 - Evaluate the hypercube layout of the generic object
* [corectl object ls](/libraries-and-tools/corectl-object-ls)	 - Print a list of all generic objects in the current app
* [corectl object properties](/libraries-and-tools/corectl-object-properties)	 - Print the properties of the generic object
* [corectl object rm](/libraries-and-tools/corectl-object-rm)	 - Remove one or many generic objects in the current app
* [corectl object set](/libraries-and-tools/corectl-object-set)	 - Set or update the objects in the current app

