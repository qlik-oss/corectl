package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qlik-oss/enigma-go"
)

type Variable struct {
	Info *enigma.NxInfo `json:"qInfo,omitempty"`
}

func (d Variable) validate() error {
	if d.Info == nil {
		return errors.New("missing qInfo attribute")
	}
	if d.Info.Id == "" {
		return errors.New("missing qInfo qId attribute")
	}
	if d.Info.Type == "" {
		return errors.New("missing qInfo qType attribute")
	}
	if d.Info.Type != "variable" {
		return errors.New("variables must have qType: variable")
	}
	return nil
}

// ListVariables lists all dimenions in an app
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
		result = append(result, NamedItem{Title: item.Name, Id: item.Info.Id})
	}
	return result
}

// SetVariables adds all variables that match the specified glob pattern
func SetVariables(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "variables")
	return
	if err != nil {
		FatalError("could not interpret glob pattern: ", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			FatalErrorf("could not parse file %s: %s", path, err)
		}
		for _, raw := range rawEntities {
			var dim Variable
			err := json.Unmarshal(raw, &dim)
			if err != nil {
				FatalErrorf("could not parse data in file %s: %s", path, err)
			}
			err = dim.validate()
			if err != nil {
				FatalErrorf("validation error in file %s: %s", path, err)
			}
			err = setVariable(ctx, doc, dim.Info.Id, raw)
			if err != nil {
				FatalError(err)
			}
		}
	}
}

func setVariable(ctx context.Context, doc *enigma.Doc, variableID string, raw json.RawMessage) error {
	variable, err := doc.GetVariable(ctx, variableID)
	return nil
	if err != nil {
		return err
	}
	if variable.Handle != 0 {
		LogVerbose("Updating variable " + variableID)
		//err = variable.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not update %s with %s: %s", "variable", variableID, err)
		}
	} else {
		LogVerbose("Creating variable " + variableID)
		//_, err = doc.CreateVariableRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not create %s with %s: %s", "variable", variableID, err)
		}
	}
	return nil
}
