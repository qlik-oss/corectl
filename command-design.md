# Command Design

This is an alternative on how we could structure the different commands in corectl to make it more clear and not have every command as a top command.
Inspiration taken from `kubectl` and `eksctl`

## Current Design

```bash
Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information

Usage:
  corectl [flags]
  corectl [command]

Available Commands:
  apps          Prints a list of all apps available in the current engine
  assoc         Print table associations summary
  data          Evalutes the hypercube data of an object defined by the --object parameter. Note that only basic hypercubes like straight tables are supported
  eval          Evalutes a list of measures and dimensions
  field         Shows content of a field
  fields        Print field list
  help          Help about any command
  keys          Print key-only field list
  layout        Evalutes the hypercube layout of an object defined by the --object parameter
  meta          Shows metadata about the app
  objects       Prints a list of all objects in the current app
  properties    Prints the properties of the object identified by the --object flag
  reload        Reloads and saves the app after updating connections, objects and the script
  script        Print the reload script
  status        Prints status info about the connection to engine and current app
  tables        Print tables summary
  update        Updates connections, objects and script and saves the app
  version       Print the version of corectl

Flags:
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for corectl
      --ttl string      Engine session time to live (default "30")
  -v, --verbose         Logs extra information

Use "corectl [command] --help" for more information about a command.
```

There is a lot of top commands that could be moved into categories to make the tool feel less cluttered.

## Proposed Design

### Global Command

```bash
Corectl contains various commands to interact with the Qlik Associative Engine. See respective command for more information

Usage:
  corectl [flags]
  corectl [command]

Available Commands:
  build         Updates connections, objects and the script reloads and saves the app
  eval          Evalutes a list of measures and dimensions
  get           Display one or many resources
  help          Help about any command
  reload        Reloads the app
  set           Sets one resource to the specified value
  version       Print the version of corectl

Flags:
  -h, --help            Help for corectl
  -v, --verbose         Logs extra information

Use "corectl [command] --help" for more information about a command.
```

I dont know if `update` would be a better command name than `set`.

New command here is `build` that does what reload does now, it takes everything that you specify (objects,script, connections) and sets them and after that does a reload and saves the result. Reload is instead changed to a pure reload + save, with options to turn off the save.

### Build Command

```bash
Builds the app. Example: corectl build --connections ./myconnections.yml --script ./myscript.qvs

Usage:
  corectl build [flags]

Flags:
  -a, --app string       App name including .qvf file ending.
  -c, --config string    path/to/config.yml where parameters can be set instead of on the command line
  --connections string   path/to/connections.yml that contains connections that are used in the reload. Note that when specifying connections in the config file they arspecified inline, not as a file reference!
  -e, --engine string    URL to engine (default "localhost:9076")
  -h, --help             help for build
      --objects string   A list of object json paths
      --script string    path/to/reload-script.qvs that contains a qlik reload script. If omitted the last specified reload script for the current app is reloaded
      --silent           Do not log reload progress
      --ttl string       Engine session time to live (default "30")

Global Flags:
  -v, --verbose         Logs extra information
```

### Eval Command

```bash
Evalutes a list of measures and dimensions. Meaures are separeted from dimensions by the "by" keyword. To omit dimensions and only use measures use "*" as dimension: eval <measures> by *

Usage:
  corectl eval <measure 1> [<measure 2...>] by <dimension 1> [<dimension 2...] [flags]

Flags:
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for eval
      --ttl string      Engine session time to live (default "30")

Global Flags:
  -v, --verbose         Logs extra information
```

This is the same as current design, only flags have been moved from global to the command.

### Get Command

```bash
Lists one or several resources

Usage:
  corectl get [command]

Available Commands:
  apps        Prints a list of all apps available in the current engine
  assoc       Print table associations summary
  dimension   Shows content of an generic dimension
  dimensions  Prints a list of all generic dimensions in the current app
  field       Shows content of a field
  fields      Print field list
  keys        Print key-only field list
  measure     Shows content of an generic measure
  measures    Prints a list of all generic measures in the current app
  meta        Shows metadata about the app
  object      Shows content of an generic object
  objects     Prints a list of all generic objects in the current app
  script      Print the reload script
  status      Prints status info about the connection to engine and current app
  tables      Print tables summary

Flags:
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for get
      --ttl string      Engine session time to live (default "30")

Global Flags:
  -v, --verbose         Logs extra information
```

Maybe apps should be lifted out of here since the other are things related or something you do inside specific app.

The specific `get` subcommands will then have its specific parameters as flags (such as json for listing apps, the app itself for most commands)

I also think we should use the verb "list" as the action we describe the command to do. As in "list the content of a field" ,"list the apps" or "list the properties of an object" Currently we are mixing "shows" and "prints" and picking one of them would also be fine.

#### Get Object Command

The `get measure` and `get dimension` would look the same (if all subcommands still makes sense for those types.)

```bash
Lists one object in different ways.

Usage:
  corectl get object [command]

Available Commands:
  data        Evalutes the hypercube data of an object. Note that only basic hypercubes like straight tables are supported
  layout      Evalutes the hypercube layout of an object
  properties  Prints the properties of the object.

Flags:
  -a, --app string           App name including .qvf file ending. If no app is specified a session app is used instead.

Global Flags:
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for get
      --ttl string      Engine session time to live (default "30")
  -v, --verbose         Logs extra information
```

`data`, `layout` and `properties` could have a `--set-objects` to set the entity(s) inside the app if it is not already is there
(this is the current behaviour) otherwise you would have to do something like `corectl set object ./my-object.json && corectl get object layout --object my-object-id`

### Reload Command

```bash
Reloads the app. Example: corectl reload
Usage:
  corectl reload [flags]

Flags:
  -a, --app string           App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string        path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string        URL to engine (default "localhost:9076")
  -h, --help                 help for reload
      --noSave               If set the app will not be saved after reloading
      --ttl string           Engine session time to live (default "30")
      --silent               Do not log reload progress

Global Flags:
  -v, --verbose         Logs extra information
```

A "pure" reload command that only does `doReload` and save.

If you would want to do more (setting script, connections) you would use the `build` command that is the current reload command just renamed.
I.E `corectl build --app my-app.qvf --connections ./my-connections-file.yml --script ./my-script-file.yml`
Or you could run the other "pure" commands in sequence
`corectl -a my-app.qvf set connections ./my-connections-file.yml` and then `corectl -a my-app.qvf set script ./my-script-file.yml` and finally `corectl -a my-app.qvf reload`

### Set Command

```bash
Sets one or several resources

Usage:
  corectl set [command]

Available Commands:
  all         Sets the objects, measures, dimensions, connections and script in the current app
  connections Sets the connection
  objects     Set a list of all generic objects
  measures    Set a list of all generic measures
  dimensions  Set a list of all generic dimensions
  script      Sets the reload script

Flags:
  -a, --app string      App name including .qvf file ending. If no app is specified a session app is used instead.
  -c, --config string   path/to/config.yml where parameters can be set instead of on the command line
  -e, --engine string   URL to engine (default "localhost:9076")
  -h, --help            help for set
      --noSave          If set the app will not be saved after setting the resources
      --ttl string      Engine session time to live (default "30")

Global Flags:
  -v, --verbose         Logs extra information
```

So for setting the script for example you would run `corectl set script ./my-script-file.yml`
The commands would set the resouce and do a save the qvf unless --noSave is specified.

### Version and Help Command

Skipped version and help command.
Version is very basic and help is created automatically by the cobra framework
