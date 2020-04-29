## corectl raw

Send Http API Request to Qlik Sense Cloud editions

### Synopsis

Send Http API Request to Qlik Sense Cloud editions. Query parameters are specified using the --query flag, a body can be specified using one of the body flags (body, body-file or body-values)

```
corectl raw <get/put/patch/post/delete> v1/url [flags]
```

### Examples

```
corectl raw get v1/items --query name=ImportantApp
```

### Options

```
      --body string                  The content of the body as a string
      --body-file string             A file path pointing to a file containing the body of the http request
      --body-values stringToString   A set of key=value pairs that well be compiled into a json object. A dot (.) inside the key is used to traverse into nested objects. The key suffixes :bool or :number can be appended to the key to inject the value into the json structure as boolean or number respectively. (default [])
  -h, --help                         help for raw
      --output-file string           A file path pointing to where the response body shoule be written
      --query stringToString         Query parameters specified as key=value pairs separated by comma (default [])
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

* [corectl](corectl.md)	 - 

