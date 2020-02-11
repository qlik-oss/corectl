package internal

import (
	"encoding/json"
	"github.com/qlik-oss/enigma-go"
	"io/ioutil"
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
