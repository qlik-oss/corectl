package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

type (
	// ParsedEntityListData struct
	ParsedEntityListData struct {
		Title string `json:"title"`
	}

	// NamedItem struct
	NamedItem struct {
		ID    string `json:"qId"`
		Title string `json:"title"`
	}

	// NamedItemWithType struct
	NamedItemWithType struct {
		ID    string `json:"qId"`
		Type  string `json:"qType,omitempty"`
		Title string `json:"title"`
	}

	// PropsWithTitle struct
	PropsWithTitle struct {
		*enigma.GenericObjectProperties
		Title string `json:"title"`
	}
)

// If commandLineGlobPattern is set the paths will be from that glob pattern
// otherwise glob patterns will be from the config file using the configEntityParam
func getEntityPaths(commandLineGlobPattern string, configEntityParam string) ([]string, error) {
	var paths []string
	var err error
	if commandLineGlobPattern != "" {
		paths, err = filepath.Glob(commandLineGlobPattern)
		if err != nil {
			return paths, err
		}
		if len(paths) == 0 {
			log.Warnf("No '%s' found for pattern %s\n", configEntityParam, commandLineGlobPattern)
		}
	} else {
		if ConfigDir == "" {
			return paths, nil
		}
		currentWorkingDir, _ := os.Getwd()
		defer os.Chdir(currentWorkingDir)
		os.Chdir(ConfigDir)

		globPatterns := viper.GetStringSlice(configEntityParam)
		var pathMatches []string
		for _, pattern := range globPatterns {
			pathMatches, err = filepath.Glob(pattern)
			if err != nil {
				log.Fatalf("could not interpret glob pattern '%s': %s\n", pattern, err)
			} else if len(pathMatches) == 0 {
				log.Warnf("No '%s' found for pattern %s\n", configEntityParam, pattern)
			} else {
				paths = append(paths, pathMatches...)
			}
		}
		for i, path := range paths {
			paths[i] = ConfigDir + "/" + path
		}
	}
	return paths, nil
}

// Try to interpret the file contents as a slice of json objects,
// otherwise try to interpret it as an json object and put it into a slice
func parseEntityFile(path string) (entities []json.RawMessage, err error) {
	entities = []json.RawMessage{}
	var entity json.RawMessage
	err = nil
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, &entities)
	if err != nil {
		entities = []json.RawMessage{}
		err = json.Unmarshal(content, &entity)
		if err == nil {
			entities = append(entities, entity)
		}
	}
	return
}
