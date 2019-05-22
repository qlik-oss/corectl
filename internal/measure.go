package internal

import (
	"fmt"
	"errors"
	"context"
	"encoding/json"

	"github.com/qlik-oss/enigma-go"
)

type Measure struct {
	Info	*enigma.NxInfo	`json:"qInfo,omitempty"`
}

func (m Measure) validate() error {
	if m.Info == nil {
		return errors.New("missing qInfo attribute")
	}
	if m.Info.Id == "" {
		return errors.New("missing qInfo qId attribute")
	}
	if m.Info.Type == "" {
		return errors.New("missing qInfo qType attribute")
	}
	if m.Info.Type != "measure" {
		return errors.New("measures must have qType: measure")
	}
	return nil
}

func ListMeasures(ctx context.Context, doc *enigma.Doc) []NamedItem {
	props := &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "corectl_entity_list",
		},
		MeasureListDef: &enigma.MeasureListDef{
			Type: "measure",
			Data: json.RawMessage(`{
				"id":"/qInfo/qId",
				"title":"/qMetaDef/title"
			}`),
		},
	}
	sessionObject, _ := doc.CreateSessionObject(ctx, props)
	defer doc.DestroySessionObject(ctx, sessionObject.GenericId)
	layout, _ := sessionObject.GetLayout(ctx)
	result := []NamedItem{}
	for _, item := range layout.MeasureList.Items {
		parsedRawData := &ParsedEntityListData{}
		json.Unmarshal(item.Data, parsedRawData)
		result = append(result, NamedItem{Title: parsedRawData.Title, Id: item.Info.Id})
	}
	return result
}

func SetMeasures(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "measures")
	if err != nil {
		FatalError("Failed to interpret glob pattern:", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			FatalError(fmt.Errorf("Failed to parse file %s: %s", path, err))
		}
		for _, raw := range rawEntities {
			var measure Measure
			err := json.Unmarshal(raw, &measure)
			if err != nil {
				FatalError(fmt.Errorf("Failed to parse data in file %s: %s", path, err))
			}
			err = measure.validate()
			if err != nil {
				FatalError(fmt.Errorf("Validation error in file %s: %s", path, err))
			}
			err = setMeasure(ctx, doc, measure.Info.Id, raw)
			if err != nil {
				FatalError("Error while creating/updating measure: ", err)
			}
		}
	}
}

func setMeasure(ctx context.Context, doc *enigma.Doc, measureID string, raw json.RawMessage) error {
	measure, err := doc.GetMeasure(ctx, measureID)
	if err != nil {
		return err
	}
	if measure.Handle != 0 {
		LogVerbose("Updating measure " + measureID)
		err = measure.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("failed to update %s with %s: %s", "measure", measureID, err)
		}
	} else {
		LogVerbose("Creating measure " + measureID)
		_, err = doc.CreateMeasureRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("failed to create %s with %s: %s", "measure", measureID, err)
		}
	}
	return nil
}
