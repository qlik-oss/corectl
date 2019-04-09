## corectl



### Synopsis

corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information

```
corectl [flags]
```

### Options

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
  -h, --help                     help for corectl
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl app](corectl_app.md)	 - Explore and manage the app
* [corectl assoc](corectl_assoc.md)	 - Print table associations summary
* [corectl build](corectl_build.md)	 - Reloads and saves the app after updating connections, dimensions, measures, objects and the script
* [corectl catwalk](corectl_catwalk.md)	 - Opens the specified app in catwalk
* [corectl completion](corectl_completion.md)	 - Generates auto completion scripts
* [corectl connection](corectl_connection.md)	 - Explore and manage connections
* [corectl dimension](corectl_dimension.md)	 - Explore and manage dimensions
* [corectl eval](corectl_eval.md)	 - Evaluates a list of measures and dimensions
* [corectl field](corectl_field.md)	 - Shows content of a field
* [corectl fields](corectl_fields.md)	 - Print field list
* [corectl keys](corectl_keys.md)	 - Print key-only field list
* [corectl measure](corectl_measure.md)	 - Explore and manage measures
* [corectl meta](corectl_meta.md)	 - Shows metadata about the app
* [corectl object](corectl_object.md)	 - Explore and manage generic objects
* [corectl reload](corectl_reload.md)	 - Reloads and saves the app.
* [corectl script](corectl_script.md)	 - Explore and manage the script
* [corectl status](corectl_status.md)	 - Prints status info about the connection to the engine and current app
* [corectl tables](corectl_tables.md)	 - Print tables summary
* [corectl version](corectl_version.md)	 - Print the version of corectl

