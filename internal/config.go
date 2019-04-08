package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/spf13/viper"
)

// ConnectionConfigEntry defines the content of a connection in either the project config yml file or a connections yml file.
type ConnectionConfigEntry struct {
	Type             string
	Username         string
	Password         string
	ConnectionString string
	Settings         map[string]string
}

// ConnectionsConfigFile defines the content of a connections yml file.
type ConnectionsConfigFile struct {
	Connections map[string]ConnectionConfigEntry
}

// FileWithReferenceToConfigFile defines a config file with a path reference to an external connections.yml file.
type FileWithReferenceToConfigFile struct {
	Connections string
}

func ResolveConnectionsFileReferenceInConfigFile(projectPath string) string {
	var configFileWithFileReference FileWithReferenceToConfigFile
	source, err := ioutil.ReadFile(projectPath)
	if err != nil {
		return projectPath
	}
	err = yaml.Unmarshal(source, &configFileWithFileReference)
	if err != nil {
		return projectPath
	}
	if configFileWithFileReference.Connections == "" {
		return projectPath
	}
	return RelativeToProject(projectPath, configFileWithFileReference.Connections)

}

// ReadConnectionsFile reads the connections config file from the supplied path.
func ReadConnectionsFile(path string) ConnectionsConfigFile {
	var config ConnectionsConfigFile
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
		explicitConfigFile = strings.TrimSpace(explicitConfigFile)
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
		LogVerbose("Using config file: " + configFile)
	} else {
		LogVerbose("No config file")
	}
}

// validateProps reads a config file by 
func validateProps(configPath string) {
	source, err := ioutil.ReadFile(configPath)
	if err != nil {
		FatalError("Could not find config file:", configPath)
	}
	validProps := map[string]struct{}{ // This "set" contains the valid property names
		"app":{}, "engine":{}, "measures":{}, "script": {},
		"dimensions":{}, "objects":{}, "connections":{},
		"headers": {}, "verbose": {}, "traffic": {},
		"no-data": {}, "bash": {},
	}
	configProps := map[string]interface{}{}
	err = yaml.Unmarshal(source, &configProps)
	if err != nil {
		FatalError(err)
	}
	invalidProps := []string{}
	for key, _ := range configProps {
		if _, ok := validProps[key]; !ok {
			invalidProps = append(invalidProps, key)
		}
	}
	if len(invalidProps) > 0 {
		errorMessage := fmt.Sprintf("Found invalid config properties: %v", invalidProps)
		FatalError(errorMessage)
	}
}

// findConfigFile finds a file with the given fileName with yml or yaml extension.
func findConfigFile(fileName string) string {
	configFile := ""
	if _, err := os.Stat(fileName + ".yml"); !os.IsNotExist(err) {
		configFile = fileName + ".yml"
	} else if _, err := os.Stat(fileName + ".yaml"); !os.IsNotExist(err) {
		configFile = fileName + ".yaml"
	}
	return configFile
}
