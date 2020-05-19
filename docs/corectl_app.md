---
title: "corectl app"
description: "corectl app"
categories: Libraries & Tools
type: Tools
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl app

Explore and manage apps

### Synopsis

Explore and manage apps

### Options

```
  -h, --help   help for app
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
* [corectl app import](/commands/corectl_app_import)	 - Import the specified app into the engine, returns the ID of the created app
* [corectl app ls](/commands/corectl_app_ls)	 - Print a list of all apps available in the current engine
* [corectl app rm](/commands/corectl_app_rm)	 - Remove the specified app

