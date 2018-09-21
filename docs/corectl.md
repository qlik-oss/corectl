## corectl



### Synopsis

Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information

```
corectl [flags]
```

### Options

```
  -a, --app string              App name including .qvf file ending
  -c, --config string           path/to/config.yml where default parameters can be set
  -e, --engine string           URL to engine
      --engine-headers string   HTTP headers to send to the engine (default "30")
  -h, --help                    help for corectl
      --ttl string              Engine session time to live (default "30")
  -v, --verbose                 Logs extra information
```

### SEE ALSO

* [corectl assoc](corectl_assoc.md)	 - Print table associations summary
* [corectl eval](corectl_eval.md)	 - Evalutes a hypercube
* [corectl field](corectl_field.md)	 - Shows content of a field
* [corectl fields](corectl_fields.md)	 - Print field list
* [corectl keys](corectl_keys.md)	 - Print key-only field list
* [corectl meta](corectl_meta.md)	 - Shows metadata about the app
* [corectl reload](corectl_reload.md)	 - Reloads the app
* [corectl script](corectl_script.md)	 - Print the reload script
* [corectl status](corectl_status.md)	 - Prints status info about the connection to engine and current app
* [corectl tables](corectl_tables.md)	 - Print tables summary

