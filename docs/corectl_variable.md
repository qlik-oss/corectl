---
title: "corectl variable"
description: "corectl variable"
categories: Libraries & Tools
type: Tools
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl variable

Explore and manage variables

### Synopsis

Explore and manage variables

### Options

```
  -h, --help   help for variable
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

* [corectl](/commands/corectl)	 - 
* [corectl variable layout](/commands/corectl_variable_layout)	 - Evaluate the layout of an generic variable
* [corectl variable ls](/commands/corectl_variable_ls)	 - Print a list of all generic variables in the current app
* [corectl variable properties](/commands/corectl_variable_properties)	 - Print the properties of the generic variable
* [corectl variable rm](/commands/corectl_variable_rm)	 - Remove one or many variables in the current app
* [corectl variable set](/commands/corectl_variable_set)	 - Set or update the variables in the current app

