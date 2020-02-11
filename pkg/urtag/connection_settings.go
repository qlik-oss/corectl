package urtag

import (
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/pkg/huggorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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
	Connections *map[string]ConnectionConfigEntry
}

// ReadConnectionsFile reads the connections config file from the supplied path.
func ReadConnectionsFile(path string) *ConnectionsConfig {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("could not find connections config file '%s'\n", path)
	}
	tempConfig := map[interface{}]interface{}{}
	err = yaml.Unmarshal(source, &tempConfig)
	if err != nil {
		log.Fatalf("invalid syntax in connections config file '%s': %s\n", path, err)
	}
	err = huggorm.SubEnvVars(&tempConfig)
	if err != nil {
		log.Fatalf("bad substitution in '%s': %s\n", path, err)
	}
	config := &ConnectionsConfig{}
	if strConfig, err := huggorm.ConvertMap(tempConfig); err == nil {
		huggorm.ReMarshal(strConfig, config)
	} else {
		log.Fatalf("could not parse connections config file '%s': %s\n", path, err)
	}
	return config
}
