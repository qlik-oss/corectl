package internal

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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
		fmt.Println("Could not find connections file:", path)
		os.Exit(1)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		FatalError(err)
	}
	return config
}
