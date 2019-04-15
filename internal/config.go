package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"bytes"
	"errors"

	"github.com/spf13/viper"
	leven "github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/yaml.v2"
)

// ConfigDir represents the directory of the config file used.
var ConfigDir string

// validProps is the set of valid config properties.
var validProps map[string]struct{} = map[string]struct{}{}

// ConnectionConfigEntry defines the content of a connection in either the project config yml file or a connections yml file.
type ConnectionConfigEntry struct {
	Type             string
	Username         string
	Password         string
	ConnectionString string
	Settings         map[string]string
}

// ConnectionsConfig represents how the connections are configured. 
type ConnectionsConfig struct {
	Connections map[string]ConnectionConfigEntry
}

// GetConnectionsConfig returns a the current connections configuration.
func GetConnectionsConfig() ConnectionsConfig {
	var config ConnectionsConfig
	conn := viper.Get("connections")
	switch conn.(type) {
	case string:
		connFile := RelativeToProject(conn.(string))
		config = ReadConnectionsFile(connFile)
	case map[string]interface{}:
		connMap := conn.(map[string]interface{})
		err := reMarshal(connMap, &config.Connections)
		if err != nil {
			FatalError(err)
		}
	}
	return config
}

// reMarshal takes a map and tries to fit it to a struct
func reMarshal(m map[string]interface{}, ref interface{}) error {
	bytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, ref)
	if err != nil {
		return err
	}
	return nil
}

// convertMap turns {} -> {} map into string -> {} map
// returns error if non-string was present in input map
func convertMap(m map[interface{}]interface{}) (map[string]interface{}, error) {
	strMap := map[string]interface{}{}
	for k, v := range(m) {
		if s, ok := k.(string); ok {
			strMap[s] = v
		} else {
			return strMap, errors.New("Non-string found in map")
		}
	}
	return strMap, nil
}

// ReadConnectionsFile reads the connections config file from the supplied path.
func ReadConnectionsFile(path string) ConnectionsConfig {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		FatalError("Could not find connections file:", path)
	}
	var config ConnectionsConfig
	tempConfig := map[interface{}]interface{}{}
	err = yaml.Unmarshal(source, &tempConfig)
	if err != nil {
		FatalError(err)
	}
	subEnvVars(&tempConfig)
	if strConfig, err := convertMap(tempConfig); err == nil {
		reMarshal(strConfig, &config)
	} else {
		FatalError("Found non-string property in:", path)
	}
	return config
}

// ReadConfigFile checks that the config file does not contain any unknown properties
// and then, if the config is valid, reads it.
func ReadConfigFile(explicitConfigFile string) {
	configFile := "" // Just for logging
	if explicitConfigFile != "" {
		explicitConfigFile, err := filepath.Abs(strings.TrimSpace(explicitConfigFile))
		if err != nil {
			FatalError(err)
		}
		setConfigFile(explicitConfigFile)
		configFile = explicitConfigFile
	} else {
		configFile = findConfigFile("corectl") // name of config file (without extension)
		if configFile != "" {
			setConfigFile(configFile)
		}
	}
	QliVerbose = viper.GetBool("verbose")
	LogTraffic = viper.GetBool("traffic")
	if configFile != "" {
		ConfigDir = filepath.Dir(configFile)
		LogVerbose("Using config file: " + configFile)
	} else {
		LogVerbose("No config file specified, using default values.")
	}
}

// AddValidProp adds the given property to the set of valid properties.
func AddValidProp(propName string) {
	validProps[propName] = struct{}{}
}

// setConfigFile reads in a config file and processes it before providing viper with it.
func setConfigFile(configPath string) {
	source, err := ioutil.ReadFile(configPath)
	if err != nil {
		FatalError("Could not find config file:", configPath)
	}
	// Using {} -> {} map to allow the recursive function subEnvVars to be less complex
	// However, this make validateProps a tiny bit more complex
	config := map[interface{}]interface{}{}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		FatalError(err)
	}
	validateProps(config, configPath)
	subEnvVars(&config)
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		FatalError(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(configBytes))
	if err != nil {
		FatalError(err)
	}
}

// findConfigFile finds a file with the given fileName with yml or yaml extension.
// Returns absolute path
func findConfigFile(fileName string) string {
	configFile := ""
	if _, err := os.Stat(fileName + ".yml"); !os.IsNotExist(err) {
		configFile = fileName + ".yml"
	} else if _, err := os.Stat(fileName + ".yaml"); !os.IsNotExist(err) {
		configFile = fileName + ".yaml"
	}
	if configFile != "" {
		absConfig, err := filepath.Abs(configFile) // Convert to abs path
		if err != nil {
			FatalError(err)
		}
		configFile = absConfig
	}
	return configFile
}

// validateProps checks if there are unknown properties in the config
// configPath is passed for error logging purposes.
func validateProps(config map[interface{}]interface{}, configPath string) {
	invalidProps := []string{}
	suggestions := map[string]string{}
	for key := range config {
		keyString, ok := key.(string)
		if !ok {
			// If there is a non-string in the yaml, this will surely be an invalid props
			keyString = fmt.Sprint(key)
		}
		if _, ok := validProps[keyString]; !ok {
			if suggestion := getSuggestion(keyString, validProps); suggestion != "" {
				suggestions[keyString] = suggestion
			} else {
				invalidProps = append(invalidProps, fmt.Sprintf("'%s'", keyString)) // For pretty printing
			}
		}
	}
	if len(invalidProps)+len(suggestions) > 0 {
		errorMessage := []string{}
		errorMessage = append(errorMessage,
			fmt.Sprintf("corectl found invalid properties when validating the config file '%s'.", configPath))
		for key, value := range suggestions {
			errorMessage = append(errorMessage, fmt.Sprintf("  '%s': did you mean '%s'?", key, value))
		}
		if len(invalidProps) > 0 {
			prepend := "M" // Capitalize M if there were no suggestions
			if len(suggestions) > 0 {
				prepend = "Also, m" // Add also if there were suggestions
			}
			errorMessage = append(errorMessage,
				fmt.Sprintf("%systerious properties: %s", prepend, strings.Join(invalidProps, ", ")))
		}
		FatalError(strings.Join(errorMessage, "\n"))
	}
}

// subEnvVars substitutes all the environment variables with their actual values in
// a map[string]interface{}, typically the unmarshallad yaml. (recursively)
func subEnvVars(m *map[interface{}]interface{}) {
	for k, v := range *m {
		switch v.(type) {
			case string:
				envVar := v.(string)
				if strings.HasPrefix(envVar, "${") && strings.HasSuffix(envVar, "}") {
					envVar = strings.TrimSuffix(strings.TrimPrefix(envVar, "${"), "}")
					if val := os.Getenv(envVar); val != "" {
						(*m)[k] = val
					} else {
						FatalError(fmt.Sprintf("Environment variable '%s' not found.", envVar))
					}
				}
			case map[interface{}]interface{}:
				m2 := v.(map[interface{}]interface{})
				subEnvVars(&m2)
		}
	}
}

// getSuggestion finds the best matching property within the specified Levenshtein distance limit
func getSuggestion(word string, validProps map[string]struct{}) string {
	op := leven.DefaultOptions // Default is cost 1 for del & ins, and 2 for substitution
	limit := 4
	min, suggestion := limit, ""
	for key := range validProps {
		dist := leven.DistanceForStrings([]rune(word), []rune(key), op)
		if dist < min {
			min = dist
			suggestion = key
		}
	}
	return suggestion
}
