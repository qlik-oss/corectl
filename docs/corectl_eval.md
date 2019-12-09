## corectl eval

Evaluate a list of measures and dimensions

### Synopsis

Evaluate a list of measures and dimensions. To evaluate a measure for a specific dimension use the <measure> by <dimension> notation. If dimensions are omitted then the eval will be evaluated over all dimensions.

```
corectl eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...] [flags]
```

### Examples

```
corectl eval "Count(a)" // returns the number of values in field "a"
corectl eval "1+1" // returns the calculated value for 1+1
corectl eval "Avg(Sales)" by "Region" // returns the average of measure "Sales" for dimension "Region"
corectl eval by "Region" // Returns the values for dimension "Region"
```

### Options

```
  -h, --help   help for eval
```

### Options inherited from parent commands

```
  -a, --app string               Name or identifier of the app
      --certificates string      path/to/folder containing client.pem, client_key.pem and root.pem certificates
  -c, --config string            path/to/config.yml where parameters can be set instead of on the command line
      --context string           Name of the context used when connecting to Qlik Associative Engine
  -e, --engine string            URL to the Qlik Associative Engine (default "localhost:9076")
      --headers stringToString   Http headers to use when connecting to Qlik Associative Engine (default [])
      --insecure                 Enabling insecure will make it possible to connect using self signed certificates
      --json                     Returns output in JSON format if possible, disables verbose and traffic output
      --no-data                  Open app without data
  -t, --traffic                  Log JSON websocket traffic to stdout
      --ttl string               Qlik Associative Engine session time to live in seconds (default "0")
  -v, --verbose                  Log extra information
```

### SEE ALSO

* [corectl](corectl.md)	 - 

