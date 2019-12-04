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
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
  -h, --help                     help for corectl
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl app](corectl_app.md)	 - Explore and manage apps
* [corectl assoc](corectl_assoc.md)	 - Print table associations
* [corectl bookmark](corectl_bookmark.md)	 - Explore and manage bookmarks
* [corectl build](corectl_build.md)	 - Reload and save the app after updating connections, dimensions, measures, objects and the script
* [corectl catwalk](corectl_catwalk.md)	 - Open the specified app in catwalk
* [corectl completion](corectl_completion.md)	 - Generate auto completion scripts
* [corectl connection](corectl_connection.md)	 - Explore and manage connections
* [corectl context](corectl_context.md)	 - Create, update and use contexts
* [corectl dimension](corectl_dimension.md)	 - Explore and manage dimensions
* [corectl eval](corectl_eval.md)	 - Evaluate a list of measures and dimensions
* [corectl fields](corectl_fields.md)	 - Print field list
* [corectl keys](corectl_keys.md)	 - Print key-only field list
* [corectl measure](corectl_measure.md)	 - Explore and manage measures
* [corectl meta](corectl_meta.md)	 - Print tables, fields and associations
* [corectl object](corectl_object.md)	 - Explore and manage generic objects
* [corectl reload](corectl_reload.md)	 - Reload and save the app
* [corectl script](corectl_script.md)	 - Explore and manage the script
* [corectl state](corectl_state.md)	 - Explore and manage alternate states
* [corectl status](corectl_status.md)	 - Print status info about the connection to the engine and current app
* [corectl tables](corectl_tables.md)	 - Print tables
* [corectl unbuild](corectl_unbuild.md)	 - Split up an existing app into separate json and yaml files
* [corectl values](corectl_values.md)	 - Print the top values of a field
* [corectl variable](corectl_variable.md)	 - Explore and manage variables
* [corectl version](corectl_version.md)	 - Print the version of corectl

