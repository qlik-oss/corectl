package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ConnectionConfigEntry struct {
	Type     string
	Username string
	Password string
	Path     string
	Settings map[string]string
}
type ConnectionsConfigFile struct {
	Connections map[string]ConnectionConfigEntry
}

func ReadConnectionsFile(path string) ConnectionsConfigFile {

	var config ConnectionsConfigFile
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not find connections file:", path)
		os.Exit(1)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	return config
}
