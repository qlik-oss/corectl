package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

type (
	ParsedEntityListData struct {
		Title string `json:"title"`
	}

	NamedItem struct {
		Id    string `json:"qId"`
		Title string `json:"title"`
	}

	NamedItemWithType struct {
		Id    string `json:"qId"`
		Type  string `json:"qType,omitempty"`
		Title string `json:"title"`
	}

	PropsWithTitle struct {
		*enigma.GenericObjectProperties
		Title string `json:"title"`
	}
)

func getEntityPaths(globPattern string, entityType string) ([]string, error) {
	paths, err := filepath.Glob(globPattern)
	if err != nil {
		return paths, err
	}
	if len(paths) == 0 {
		if ConfigDir == "" {
			return paths, nil
		}
		currentWorkingDir, _ := os.Getwd()
		defer os.Chdir(currentWorkingDir)
		os.Chdir(ConfigDir)

		globPatterns := viper.GetStringSlice(entityType)

		for _, pattern := range globPatterns {
			paths, err = filepath.Glob(pattern)
			if err != nil {
				FatalError(err)
			}
			for i, path := range paths {
				paths[i] = ConfigDir + "/" + path
			}
		}
	}
	return paths, nil
}

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
		err = json.Unmarshal(content, &entity)
		if err == nil {
			entities = append(entities, entity)
		}
	}
	return
}
