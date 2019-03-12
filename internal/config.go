package internal

import (
	"io/ioutil"

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

// ReadConnectionsFile reads the connections config file from the supplied path.
func ReadConnectionsFile(path string) ConnectionsConfigFile {

	var config ConnectionsConfigFile
	source, err := ioutil.ReadFile(path)
	if err != nil {
		Logger.Fatalf("Could not find connections file: %s", path)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		Logger.Fatal(err)
	}
	return config
}
