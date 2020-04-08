package dynconf

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

func ReadSettings(ccmd *cobra.Command) *DynSettings {
	result := readSettings(ccmd.Flags(), true)
	root := ccmd.Root()
	result.rootName = root.Use
	result.version = root.Version
	return result
}
func ReadSettingsWithoutContext(ccmd *cobra.Command) *DynSettings {
	return readSettings(ccmd.Flags(), false)
}

func readSettings(commandFlagSet *pflag.FlagSet, withContext bool) *DynSettings {
	contextName, err := commandFlagSet.GetString("context")
	if err != nil {
		panic(err)
	}
	configFileName, err := commandFlagSet.GetString("config")
	if err != nil {
		panic(err)
	}

	var configPath string
	// Using {} -> {} map to allow the recursive function subEnvVars to be less complex
	// However, this make validateProps a tiny bit more complex
	var commandLineParams = map[string]interface{}{}
	var configParams = map[string]interface{}{}
	var contextParams = map[string]interface{}{}
	var defaultParams = map[string]interface{}{}

	// Read config file if specified.
	if configFileName == "" {
		configFileName = findConfigFile("corectl")
	}

	if configFileName != "" {
		configPath = filepath.Dir(configFileName)
		source, err := ioutil.ReadFile(configFileName)
		if err != nil {
			log.Fatalf("could not find config file '%s'\n", configFileName)
		}
		tempConfig := &map[interface{}]interface{}{}
		err = yaml.Unmarshal(source, tempConfig)
		if err != nil {
			log.Fatalf("invalid syntax in config file '%s': %s\n", configFileName, err)
		}
		validateProps(*tempConfig, configFileName)
		if err = subEnvVars(tempConfig); err != nil {
			log.Fatalf("bad substitution in '%s': %s\n", configFileName, err)
		}
		if configParams, err = convertMap(*tempConfig); err != nil {
			log.Fatalf("invalid syntax in config file '%s': %s\n", configFileName, err)
		}
	}

	// Read context, if it should be included.
	if withContext {
		contextHandler := NewContextHandler()
		if contextName == "" {
			contextName = contextHandler.Current
		}
		context := contextHandler.Get(contextName)
		contextParams = map[string]interface{}(context)
	}

	// Check for overlap in config and context.
	for k := range contextParams {
		// Headers from context will be merged with that of config, so no conflict.
		if k == "headers" {
			continue
		}
		if _, ok := configParams[k]; ok {
			log.Warnf("Property '%s' exists in both current context and config, using property from config\n", k)
		}
	}

	// Read command-line parameters
	commandFlagSet.Visit(func(flag *pflag.Flag) {
		value := getFlagValue(commandFlagSet, flag)
		commandLineParams[flag.Name] = value
	})

	// Read default values only if not set in any other parameter source.
	commandFlagSet.VisitAll(func(flag *pflag.Flag) {
		_, fromCmd := commandLineParams[flag.Name]
		_, fromCfg := configParams[flag.Name]
		_, fromCtx := contextParams[flag.Name]
		if !(fromCmd || fromCfg || fromCtx) {
			value := getFlagValue(commandFlagSet, flag)
			defaultParams[flag.Name] = value
		}
	})

	return &DynSettings{
		contextName:       contextName,
		configPath:        configPath,
		configFilePath:    configFileName,
		commandLineParams: commandLineParams,
		contextParams:     contextParams,
		configParams:      configParams,
		defaultParams:     defaultParams,
	}
}

// relativeToProject transforms a path to be relative to a base path of the project file
func relativeToProject(baseDir string, path string) string {
	if baseDir != "" && !filepath.IsAbs(path) {
		fullpath := filepath.Join(baseDir, path)
		return fullpath
	}
	return path
}

func getFlagValue(flagset *pflag.FlagSet, flag *pflag.Flag) interface{} {
	var result interface{}
	switch flag.Value.Type() {
	case "string":
		result, _ = flagset.GetString(flag.Name)
	case "int":
		result, _ = flagset.GetInt(flag.Name)
	case "bool":
		result, _ = flagset.GetBool(flag.Name)
	case "stringSlice":
		result, _ = flagset.GetStringSlice(flag.Name)
	case "stringToString":
		result, _ = flagset.GetStringToString(flag.Name)
	default:
		panic("Unexpected type:" + flag.Value.Type())
	}
	return result
}

type DynSettings struct {
	rootName string
	version  string

	contextName       string
	configPath        string
	configFilePath    string
	commandLineParams map[string]interface{}
	configParams      map[string]interface{}
	contextParams     map[string]interface{}
	defaultParams     map[string]interface{}
}

func (ds *DynSettings) OverrideSetting(name string, value interface{}) {
	ds.commandLineParams[name] = value
}

func (ds *DynSettings) ConfigPath() string {
	return ds.configPath
}

func (ds *DynSettings) ConfigFilePath() string {
	return ds.configFilePath
}

func (ds *DynSettings) get(name string) interface{} {
	if val, ok := ds.commandLineParams[name]; ok {
		return val
	}
	if val, ok := ds.configParams[name]; ok {
		return val
	}
	if val, ok := ds.contextParams[name]; ok {
		return val
	}
	return ds.defaultParams[name]
}

func (ds *DynSettings) GetString(name string) string {
	switch value := ds.get(name).(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	default:
		log.Fatalf("Unexpected type of parameter: %s", name)
		return ""
	}
}
func (ds *DynSettings) GetInt(name string) int {
	switch value := ds.get(name).(type) {
	case string:
		res, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Failed to parse boolean in parameter %s, %s", name, err)
		}
		return res
	case int:
		return value
	default:
		log.Fatalf("Unexpected type of parameter %s: %T", name, value)
		return 0
	}
}

func (ds *DynSettings) GetBoolAllowNoFlag(name string) bool {
	switch value := ds.get(name).(type) {
	case string:
		res, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("Failed to parse boolean in parameter %s, %s", name, err)
		}
		return res
	case bool:
		return value
	default:
		return false
	}
}
func (ds *DynSettings) GetBool(name string) bool {
	switch value := ds.get(name).(type) {
	case string:
		res, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("Failed to parse boolean in parameter %s, %s", name, err)
		}
		return res
	case bool:
		return value
	default:

		if name == "quiet" {
			return false
		}
		//TODO Should we allow defaults for non-present flags? log.Fatalf("Unexpected type of parameter %s: %s", name, reflect.TypeOf(value).Name())

		if value == nil {
			log.Fatalf("No such flag: %s", name)
		}

		log.Fatalf("Unexpected type of parameter %s: %T", name, value)
		return false
	}
}

func (ds *DynSettings) IsDefinedOnCommandLine(name string) bool {
	_, ok := ds.commandLineParams[name]
	return ok
}

func (ds *DynSettings) GetStringArray(name string) []string {
	switch array := ds.get(name).(type) {
	case []interface{}:
		result := make([]string, len(array))
		for i, v := range array {
			str, ok := v.(string)
			if !ok {
				return []string{}
			}
			result[i] = str
		}
		return result
	case string:
		return []string{array}
	}
	return []string{}
}

func (ds *DynSettings) IsString(name string) bool {
	_, ok := ds.get(name).(string)
	return ok
}

func (ds *DynSettings) GetStringMap(name string) map[string]string {
	value := ds.get(name)
	return toStringMap(value)
}

func toStringMap(x interface{}) map[string]string {
	switch value := x.(type) {
	case map[string]string:
		return value
	case map[interface{}]interface{}:
		result := make(map[string]string)
		for i, x := range value {
			stringKey, ok1 := i.(string)
			stringValue, ok2 := x.(string)
			if ok1 && ok2 {
				result[stringKey] = stringValue
			}
		}
		return result
	default:
		if value == nil {
			return nil
		}
		log.Fatalf("Unexpected format of map: %T", value)
		return nil
	}
}

func (ds *DynSettings) GetPath(name string) string {
	path := ds.GetString(name)
	if path == "" {
		return ""
	} else if ds.IsDefinedOnCommandLine(name) {
		currentWorkingDir, _ := os.Getwd()
		return relativeToProject(currentWorkingDir, path)
	} else {
		return relativeToProject(ds.configPath, path)
	}
}

func (ds *DynSettings) GetAbsolutePath(name string) string {
	path := ds.GetPath(name)
	if path == "" {
		return ""
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Not a valid file path: %s", path)
	}
	return absPath
}

func (ds *DynSettings) GetGlobFiles(name string) []string {
	if ds.IsDefinedOnCommandLine(name) {
		commandLineGlobPattern := ds.GetString(name)
		if commandLineGlobPattern == "" {
			return []string{}
		}
		paths, err := filepath.Glob(commandLineGlobPattern)
		if err != nil {
			log.Fatalf("could not interpret glob pattern '%s': %s\n", commandLineGlobPattern, err)
		}
		return paths
	} else {
		var paths []string
		if ds.configPath == "" {
			return paths
		}
		currentWorkingDir, _ := os.Getwd()
		defer os.Chdir(currentWorkingDir)
		os.Chdir(ds.configPath)

		globPatterns := ds.GetStringArray(name) //Allow array in yaml
		for _, pattern := range globPatterns {
			pathMatches, err := filepath.Glob(pattern)
			if err != nil {
				log.Fatalf("could not interpret glob pattern '%s': %s\n", pattern, err)
			} else {
				paths = append(paths, pathMatches...)
			}
		}
		for i, path := range paths {
			paths[i] = ds.configPath + "/" + path
		}
		return paths
	}
}

// ReadCertificates reads and loads the specified certificates
func (ds *DynSettings) GetTLSConfigFromPath(certificatesPath string) *tls.Config {
	tlsClientConfig := &tls.Config{}

	// Read client and root certificates.
	certPath := ds.GetPath(certificatesPath)
	if certPath != "" {
		certFile := certPath + "/client.pem"
		keyFile := certPath + "/client_key.pem"
		caFile := certPath + "/root.pem"

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalln("could not load client certificate: ", err)
		}

		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalln("could not read root certificate: ", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup TLS cert configuration.
		tlsClientConfig.Certificates = []tls.Certificate{cert}
		tlsClientConfig.RootCAs = caCertPool
	}
	return tlsClientConfig

}

// Returns true if no value has been set in either the context,
// config file or command line for the supplied flag name
func (ds *DynSettings) IsUsingDefaultValue(name string) bool {
	_, ok := ds.defaultParams[name]
	return ok
}

// GetHeaders merges the headers from the different parameter sources in order of precedence.
// It also, converts the key to lower-case.
func (ds *DynSettings) GetHeaders() http.Header {
	result := http.Header{}
	headers := make([]map[string]string, 3)
	headers[0] = toStringMap(ds.commandLineParams["headers"])
	headers[1] = toStringMap(ds.configParams["headers"])
	headers[2] = toStringMap(ds.contextParams["headers"])
	for _, header := range headers {
		for k, v := range header {
			if result.Get(k) == "" {
				result.Add(k, v)
			}
		}
	}
	if agent := ds.GetUserAgent(); agent != "" {
		result.Set("User-Agent", agent)
	}
	return result
}

// GetUserAgent returns a string representing the User-Agent, consisting of the root name
// of the associated command, its version and the OS.
// If no rootName is set, the returned value will be empty.
func (ds *DynSettings) GetUserAgent() string {
	var agent string
	if ds.rootName == "" {
		return ""
	}
	agent = ds.rootName
	if ds.version != "" {
		agent += "/" + ds.version
	}
	agent += " (" + runtime.GOOS + ")"
	return agent
}
