package internal

import (
	"fmt"
	"errors"
	"context"
	"encoding/json"

	"github.com/qlik-oss/enigma-go"
)

type Dimension struct {
	Info	*enigma.NxInfo	`json:"qInfo,omitempty"`
}

func (d Dimension) validate() error {
	if d.Info == nil {
		return errors.New("missing qInfo attribute")
	}
	if d.Info.Id == "" {
		return errors.New("missing qInfo qId attribute")
	}
	if d.Info.Type == "" {
		return errors.New("missing qInfo qType attribute")
	}
	if d.Info.Type != "dimension" {
		return errors.New("dimensions must have qType: dimension")
	}
	return nil
}

func ListDimensions(ctx context.Context, doc *enigma.Doc) []NamedItem {
	props := &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "corectl_entity_list",
		},
		DimensionListDef: &enigma.DimensionListDef{
			Type: "dimension",
			Data: json.RawMessage(`{
				"id":"/qInfo/qId",
				"title":"/qDim/title"
			}`),
		},
	}
	sessionObject, _ := doc.CreateSessionObject(ctx, props)
	defer doc.DestroySessionObject(ctx, sessionObject.GenericId)
	layout, _ := sessionObject.GetLayout(ctx)
	result := []NamedItem{}
	for _, item := range layout.DimensionList.Items {
		parsedRawData := &ParsedEntityListData{}
		json.Unmarshal(item.Data, parsedRawData)
		result = append(result, NamedItem{Title: parsedRawData.Title, Id: item.Info.Id})
	}
	return result
}

func SetDimensions(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "dimensions")
	if err != nil {
		FatalError("Failed to interpret glob pattern:", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			FatalError(fmt.Errorf("Failed to parse file %s: %s", path, err))
		}
		for _, raw := range rawEntities {
			var dimension Dimension
			err := json.Unmarshal(raw, &dimension)
			if err != nil {
				FatalError(fmt.Errorf("Failed to parse data in file %s: %s", path, err))
			}
			err = dimension.validate()
			if err != nil {
				FatalError(fmt.Errorf("Validation error in file %s: %s", path, err))
			}
			err = setDimension(ctx, doc, dimension.Info.Id, raw)
			if err != nil {
				FatalError("Error while creating/updating dimension: ", err)
			}
		}
	}
}

func setDimension(ctx context.Context, doc *enigma.Doc, dimensionID string, raw json.RawMessage) error {
	dimension, err := doc.GetDimension(ctx, dimensionID)
	if err != nil {
		return err
	}
	if dimension.Handle != 0 {
		LogVerbose("Updating dimension " + dimensionID)
		err = dimension.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("failed to update %s with %s: %s", "dimension", dimensionID, err)
		}
	} else {
		LogVerbose("Creating dimension " + dimensionID)
		_, err = doc.CreateDimensionRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("failed to create %s with %s: %s", "dimension", dimensionID, err)
		}
	}
	return nil
}
