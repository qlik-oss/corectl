---
title: "corectl script set"
description: "corectl script set"
categories: Libraries & Tools
type: Tools
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl script set

Set the script in the current app

### Synopsis

Set the script in the current app

```
corectl script set <path-to-script-file.qvs> [flags]
```

### Examples

```
corectl script set ./my-script-file.qvs
```

### Options

```
  -h, --help      help for set
      --no-save   Do not save the app
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

* [corectl script](/commands/corectl_script)	 - Explore and manage the script

