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

* [corectl build](corectl_build.md)	 - Reloads and saves the app after updating connections, dimensions, measures, objects and the script
* [corectl catwalk](corectl_catwalk.md)	 - Opens the specified app in catwalk
* [corectl completion](corectl_completion.md)	 - Generates auto completion scripts
* [corectl eval](corectl_eval.md)	 - Evaluates a list of measures and dimensions
* [corectl get](corectl_get.md)	 - Lists one or several resources
* [corectl reload](corectl_reload.md)	 - Reloads the app.
* [corectl remove](corectl_remove.md)	 - Remove entities (connections, dimensions, measures, objects) in the app or the app itself
* [corectl set](corectl_set.md)	 - Sets one or several resources
* [corectl version](corectl_version.md)	 - Print the version of corectl

