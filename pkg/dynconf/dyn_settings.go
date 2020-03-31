package dynconf

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

func ReadSettings(ccmd *cobra.Command) *DynSettings {
	result := readSettings(ccmd.Flags(), true)
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
	var allParams = &map[interface{}]interface{}{}
	var commandLineParams = &map[interface{}]interface{}{}
	var defaultValueIsUsed = &map[interface{}]bool{}

	if configFileName == "" {
		configFileName = findConfigFile("corectl")
	}

	if configFileName != "" {
		configPath = filepath.Dir(configFileName)
		source, err := ioutil.ReadFile(configFileName)
		if err != nil {
			log.Fatalf("could not find config file '%s'\n", configFileName)
		}
		err = yaml.Unmarshal(source, allParams)
		if err != nil {
			log.Fatalf("invalid syntax in config file '%s': %s\n", configFileName, err)
		}
	}

	// Merge before validation and env substitution since it might not be needed due to context.
	if withContext {
		mergeContext(allParams, contextName)
	}

	validateProps(*allParams, configFileName)
	err = subEnvVars(allParams)
	if err != nil {
		log.Fatalf("bad substitution in '%s': %s\n", configFileName, err)
	}

	commandFlagSet.Visit(func(flag *pflag.Flag) {
		value := getFlagValue(commandFlagSet, flag)
		(*allParams)[flag.Name] = value
		(*commandLineParams)[flag.Name] = value
	})

	commandFlagSet.VisitAll(func(flag *pflag.Flag) {
		if (*allParams)[flag.Name] == nil {
			value := getFlagValue(commandFlagSet, flag)
			(*allParams)[flag.Name] = value
			(*defaultValueIsUsed)[flag.Name] = true
		} else {
			(*defaultValueIsUsed)[flag.Name] = false
		}
	})

	return &DynSettings{
		contextName:        contextName,
		configPath:         configPath,
		configFilePath:     configFileName,
		allParams:          *allParams,
		commandLineParams:  *commandLineParams,
		defaultValueIsUsed: *defaultValueIsUsed,
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
	contextName        string
	configPath         string
	configFilePath     string
	allParams          map[interface{}]interface{}
	commandLineParams  map[interface{}]interface{}
	defaultValueIsUsed map[interface{}]bool
}

func (ds *DynSettings) OverrideSetting(name string, value interface{}) {
	ds.allParams[name] = value
}
func (ds *DynSettings) ConfigPath() string {
	return ds.configPath
}
func (ds *DynSettings) ConfigFilePath() string {
	return ds.configFilePath
}

func (ds *DynSettings) GetString(name string) string {
	switch filevalue := ds.allParams[name].(type) {
	case string:
		return filevalue
	case int:
		return strconv.Itoa(filevalue)
	default:
		log.Fatalf("Unexpected type of parameter: %s", name)
		return ""
	}
}
func (ds *DynSettings) GetInt(name string) int {
	switch value := ds.allParams[name].(type) {
	case string:
		res, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Failed to parse boolean in parameter %s, %s", name, err)
		}
		return res
	case int:
		return value
	default:
		log.Fatalf("Unexpected type of parameter %s: %s", name, reflect.TypeOf(value).Name())
		return 0
	}
}

func (ds *DynSettings) GetBoolAllowNoFlag(name string) bool {
	switch value := ds.allParams[name].(type) {
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
	switch value := ds.allParams[name].(type) {
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

		log.Fatalf("Unexpected type of parameter %s: %s", name, reflect.TypeOf(value).Name())
		return false
	}
}

func (ds *DynSettings) IsDefinedOnCommandLine(name string) bool {
	return ds.commandLineParams[name] != nil
}
func (ds *DynSettings) GetStringArray(name string) []string {
	switch array := ds.allParams[name].(type) {
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
	switch ds.allParams[name].(type) {
	case string:
		return true
	default:
		return false
	}
}

func (ds *DynSettings) GetStringMap(name string) map[string]string {
	switch value := ds.allParams[name].(type) {
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
		log.Fatalf("Unexpected format of map: %s", reflect.TypeOf(value).Name())
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

// Returns true if no value has been set in either the context, config file or command line for the supplied flag name
func (ds *DynSettings) IsUsingDefaultValue(name string) bool {
	return ds.defaultValueIsUsed[name]
}
