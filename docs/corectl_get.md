## corectl get

Lists one or several resources

### Synopsis

Lists one or several resources

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -a, --app string               App name, if no app is specified a session app is used instead.
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "30")
  -v, --verbose                  Logs extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 
* [corectl get apps](corectl_get_apps.md)	 - Prints a list of all apps available in the current engine
* [corectl get assoc](corectl_get_assoc.md)	 - Print table associations summary
* [corectl get connection](corectl_get_connection.md)	 - Shows the properties for a specific connection
* [corectl get connections](corectl_get_connections.md)	 - Prints a list of all connections in the specified app
* [corectl get dimension](corectl_get_dimension.md)	 - Shows content of an generic dimension
* [corectl get dimensions](corectl_get_dimensions.md)	 - Prints a list of all generic dimensions in the current app
* [corectl get field](corectl_get_field.md)	 - Shows content of a field
* [corectl get fields](corectl_get_fields.md)	 - Print field list
* [corectl get keys](corectl_get_keys.md)	 - Print key-only field list
* [corectl get measure](corectl_get_measure.md)	 - Shows content of an generic measure
* [corectl get measures](corectl_get_measures.md)	 - Prints a list of all generic measures in the current app
* [corectl get meta](corectl_get_meta.md)	 - Shows metadata about the app
* [corectl get object](corectl_get_object.md)	 - Shows content of an generic object
* [corectl get objects](corectl_get_objects.md)	 - Prints a list of all generic objects in the current app
* [corectl get script](corectl_get_script.md)	 - Print the reload script
* [corectl get status](corectl_get_status.md)	 - Prints status info about the connection to the engine and current app
* [corectl get tables](corectl_get_tables.md)	 - Print tables summary

