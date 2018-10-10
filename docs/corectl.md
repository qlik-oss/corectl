## corectl



### Synopsis

Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information

```
corectl [flags]
```

### Options

```
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for corectl
      --ttl string      Engine session time to live (default "30")
  -v, --verbose         Logs extra information
```

### SEE ALSO

* [corectl apps](corectl_apps.md)	 - Prints a list of all apps available in the current engine
* [corectl assoc](corectl_assoc.md)	 - Print table associations summary
* [corectl data](corectl_data.md)	 - Evalutes the hypercube data of an object defined by the --object parameter. Note that only basic hypercubes like straight tables are supported
* [corectl eval](corectl_eval.md)	 - Evalutes a list of measures and dimensions
* [corectl field](corectl_field.md)	 - Shows content of a field
* [corectl fields](corectl_fields.md)	 - Print field list
* [corectl keys](corectl_keys.md)	 - Print key-only field list
* [corectl layout](corectl_layout.md)	 - Evalutes the hypercube layout of an object defined by the --object parameter
* [corectl meta](corectl_meta.md)	 - Shows metadata about the app
* [corectl objects](corectl_objects.md)	 - Prints a list of all objects in the current app
* [corectl properties](corectl_properties.md)	 - Prints the properties of the object identified by the --object flag
* [corectl reload](corectl_reload.md)	 - Reloads and saves the app after updating connections, objects and the script
* [corectl script](corectl_script.md)	 - Print the reload script
* [corectl status](corectl_status.md)	 - Prints status info about the connection to engine and current app
* [corectl tables](corectl_tables.md)	 - Print tables summary
* [corectl version](corectl_version.md)	 - Print the version of corectl
* [corectl update](corectl_update.md)	 - Updates connections, objects and script and saves the app

