package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	leven "github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/yaml.v2"
)

// ConfigDir represents the directory of the config file used.
var ConfigDir string

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

// validProps is the set of valid config properties.
var validProps map[string]struct{} = map[string]struct{}{}

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

// ReadConnectionsFile reads the connections config file from the supplied path.
func ReadConnectionsFile(path string) ConnectionsConfig {
	var config ConnectionsConfig
	source, err := ioutil.ReadFile(path)
	if err != nil {
		FatalError("Could not find connections file:", path)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		FatalError(err)
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
		validateProps(explicitConfigFile)
		viper.SetConfigFile(explicitConfigFile)
		if err := viper.ReadInConfig(); err == nil {
			configFile = explicitConfigFile
		} else {
			FatalError(err)
		}
	} else {
		configFile = findConfigFile("corectl") // name of config file (without extension)
		if configFile != "" {
			validateProps(configFile)
			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				FatalError("Failed to read config file " + configFile)
			}
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

// validateProps reads a config file by
func validateProps(configPath string) {
	source, err := ioutil.ReadFile(configPath)
	if err != nil {
		FatalError("Could not find config file:", configPath)
	}
	configProps := map[string]interface{}{}
	err = yaml.Unmarshal(source, &configProps)
	if err != nil {
		FatalError(err)
	}
	invalidProps := []string{}
	suggestions := map[string]string{}
	for key, _ := range configProps {
		if _, ok := validProps[key]; !ok {
			if suggestion := getSuggestion(key, validProps); suggestion != "" {
				suggestions[key] = suggestion
			} else {
				invalidProps = append(invalidProps, fmt.Sprintf("'%s'", key)) // For pretty printing
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

// getSuggestion finds the best matching property within the specified Levenshtein distance limit
func getSuggestion(word string, validProps map[string]struct{}) string {
	op := leven.DefaultOptions // Default is cost 1 for del & ins, and 2 for substitution
	limit := 4
	min, suggestion := limit, ""
	for key, _ := range validProps {
		dist := leven.DistanceForStrings([]rune(word), []rune(key), op)
		if dist < min {
			min = dist
			suggestion = key
		}
	}
	return suggestion
}
