---
title: "corectl measure"
description: "corectl measure"
categories: Libraries & Tools
type: Tools
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl measure

Explore and manage measures

### Synopsis

Explore and manage measures

### Options

```
  -h, --help   help for measure
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
* [corectl measure layout](/commands/corectl_measure_layout)	 - Evaluate the layout of an generic measure
* [corectl measure ls](/commands/corectl_measure_ls)	 - Print a list of all generic measures in the current app
* [corectl measure properties](/commands/corectl_measure_properties)	 - Print the properties of the generic measure
* [corectl measure rm](/commands/corectl_measure_rm)	 - Remove one or many generic measures in the current app
* [corectl measure set](/commands/corectl_measure_set)	 - Set or update the measures in the current app

