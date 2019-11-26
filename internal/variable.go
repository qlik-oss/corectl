package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// Variable qix object struct
type Variable struct {
	Info *enigma.NxInfo `json:"qInfo,omitempty"`
	Name string         `json:"qName,omitempty"`
}

func (v Variable) validate() error {
	if v.Info == nil {
		return errors.New("missing qInfo attribute")
	}
	if v.Info.Id == "" {
		return errors.New("missing qInfo qId attribute")
	}
	if v.Info.Type == "" {
		return errors.New("missing qInfo qType attribute")
	}
	if v.Name == "" {
		return errors.New("missing Name attribute")
	}
	if v.Info.Type != "variable" {
		return errors.New("variables must have qType: variable")
	}
	return nil
}

// ListVariables lists all variables in an app
func ListVariables(ctx context.Context, doc *enigma.Doc) []NamedItem {
	props := &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "corectl_entity_list",
		},
		VariableListDef: &enigma.VariableListDef{
			Type: "variable",
			Data: json.RawMessage(`{
				"id":"/qInfo/qId",
				"title":"/qMetaDef/title",
				"name":"/qMetaDef/name"
			}`),
		},
	}
	sessionObject, _ := doc.CreateSessionObject(ctx, props)
	defer doc.DestroySessionObject(ctx, sessionObject.GenericId)
	layout, _ := sessionObject.GetLayout(ctx)
	result := []NamedItem{}
	for _, item := range layout.VariableList.Items {
		result = append(result, NamedItem{Title: item.Name, ID: item.Info.Id})
	}
	return result
}

// SetVariables adds all variables that match the specified glob pattern
func SetVariables(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "variables")
	if err != nil {
		log.Fatalln("could not interpret glob pattern: ", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			log.Fatalf("could not parse file %s: %s\n", path, err)
		}
		for _, raw := range rawEntities {
			var variable Variable
			err := json.Unmarshal(raw, &variable)
			if err != nil {
				log.Fatalf("could not parse data in file %s: %s\n", path, err)
			}
			err = variable.validate()
			if err != nil {
				log.Fatalf("validation error in file %s: %s\n", path, err)
			}
			err = setVariable(ctx, doc, variable.Name, raw)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func setVariable(ctx context.Context, doc *enigma.Doc, variableName string, raw json.RawMessage) error {
	variable, err := doc.GetVariableByName(ctx, variableName)
	if err != nil {
		return err
	}
	if variable.Handle != 0 {
		log.Verboseln("Updating variable " + variableName)
		err = variable.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not update %s with %s: %s", "variable", variableName, err)
		}
	} else {
		log.Verboseln("Creating variable " + variableName)
		_, err = doc.CreateVariableExRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not create %s with %s: %s", "variable", variableName, err)
		}
	}
	return nil
}
