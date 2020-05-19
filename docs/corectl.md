---
title: "corectl"
description: "corectl"
categories: Libraries & Tools
type: Commands
tags: qlik-cli
products: Qlik Cloud, QSEoK
---
## corectl



### Synopsis

corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information

```
corectl [flags]
```

### Options

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Name of the context used when connecting to Qlik Associative Engine
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
  -h, --help                     help for corectl
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -s, --server string            URL to a Qlik Product, a local engine, cluster or sense-enterprise
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl app](/libraries-and-tools/corectl-app)	 - Explore and manage apps
* [corectl assoc](/libraries-and-tools/corectl-assoc)	 - Print table associations
* [corectl bookmark](/libraries-and-tools/corectl-bookmark)	 - Explore and manage bookmarks
* [corectl build](/libraries-and-tools/corectl-build)	 - Reload and save the app after updating connections, dimensions, measures, objects and the script
* [corectl catwalk](/libraries-and-tools/corectl-catwalk)	 - Open the specified app in catwalk
* [corectl completion](/libraries-and-tools/corectl-completion)	 - Generate auto completion scripts
* [corectl connection](/libraries-and-tools/corectl-connection)	 - Explore and manage connections
* [corectl context](/libraries-and-tools/corectl-context)	 - Create, update and use contexts
* [corectl dimension](/libraries-and-tools/corectl-dimension)	 - Explore and manage dimensions
* [corectl eval](/libraries-and-tools/corectl-eval)	 - Evaluate a list of measures and dimensions
* [corectl fields](/libraries-and-tools/corectl-fields)	 - Print field list
* [corectl keys](/libraries-and-tools/corectl-keys)	 - Print key-only field list
* [corectl measure](/libraries-and-tools/corectl-measure)	 - Explore and manage measures
* [corectl meta](/libraries-and-tools/corectl-meta)	 - Print tables, fields and associations
* [corectl object](/libraries-and-tools/corectl-object)	 - Explore and manage generic objects
* [corectl raw](/libraries-and-tools/corectl-raw)	 - Send Http API Request to Qlik Sense Cloud editions
* [corectl reload](/libraries-and-tools/corectl-reload)	 - Reload and save the app
* [corectl script](/libraries-and-tools/corectl-script)	 - Explore and manage the script
* [corectl state](/libraries-and-tools/corectl-state)	 - Explore and manage alternate states
* [corectl status](/libraries-and-tools/corectl-status)	 - Print status info about the connection to the engine and current app
* [corectl tables](/libraries-and-tools/corectl-tables)	 - Print tables
* [corectl unbuild](/libraries-and-tools/corectl-unbuild)	 - Split up an existing app into separate json and yaml files
* [corectl values](/libraries-and-tools/corectl-values)	 - Print the top values of a field
* [corectl variable](/libraries-and-tools/corectl-variable)	 - Explore and manage variables
* [corectl version](/libraries-and-tools/corectl-version)	 - Print the version of corectl

