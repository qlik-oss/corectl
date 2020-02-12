package boot

import (
	"github.com/qlik-oss/corectl/pkg/dynconf"
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
	config := &ConnectionsConfig{}
	dynconf.ReadYamlFile(path, "connections config file", config)
	return config
}
