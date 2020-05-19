---
title: "corectl"
description: "corectl"
categories: Libraries & Tools
type: Tools
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

* [corectl app](/commands/corectl_app)	 - Explore and manage apps
* [corectl assoc](/commands/corectl_assoc)	 - Print table associations
* [corectl bookmark](/commands/corectl_bookmark)	 - Explore and manage bookmarks
* [corectl build](/commands/corectl_build)	 - Reload and save the app after updating connections, dimensions, measures, objects and the script
* [corectl catwalk](/commands/corectl_catwalk)	 - Open the specified app in catwalk
* [corectl completion](/commands/corectl_completion)	 - Generate auto completion scripts
* [corectl connection](/commands/corectl_connection)	 - Explore and manage connections
* [corectl context](/commands/corectl_context)	 - Create, update and use contexts
* [corectl dimension](/commands/corectl_dimension)	 - Explore and manage dimensions
* [corectl eval](/commands/corectl_eval)	 - Evaluate a list of measures and dimensions
* [corectl fields](/commands/corectl_fields)	 - Print field list
* [corectl keys](/commands/corectl_keys)	 - Print key-only field list
* [corectl measure](/commands/corectl_measure)	 - Explore and manage measures
* [corectl meta](/commands/corectl_meta)	 - Print tables, fields and associations
* [corectl object](/commands/corectl_object)	 - Explore and manage generic objects
* [corectl raw](/commands/corectl_raw)	 - Send Http API Request to Qlik Sense Cloud editions
* [corectl reload](/commands/corectl_reload)	 - Reload and save the app
* [corectl script](/commands/corectl_script)	 - Explore and manage the script
* [corectl state](/commands/corectl_state)	 - Explore and manage alternate states
* [corectl status](/commands/corectl_status)	 - Print status info about the connection to the engine and current app
* [corectl tables](/commands/corectl_tables)	 - Print tables
* [corectl unbuild](/commands/corectl_unbuild)	 - Split up an existing app into separate json and yaml files
* [corectl values](/commands/corectl_values)	 - Print the top values of a field
* [corectl variable](/commands/corectl_variable)	 - Explore and manage variables
* [corectl version](/commands/corectl_version)	 - Print the version of corectl

