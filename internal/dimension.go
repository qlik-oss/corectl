package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

type Dimension struct {
	Info *enigma.NxInfo `json:"qInfo,omitempty"`
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

// ListDimensions lists all dimensions in an app
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

// SetDimensions adds all dimensions that match the specified glob pattern
func SetDimensions(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "dimensions")
	if err != nil {
		log.Fatalln("could not interpret glob pattern: ", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			log.Fatalf("could not parse file %s: %s\n", path, err)
		}
		for _, raw := range rawEntities {
			var dim Dimension
			err := json.Unmarshal(raw, &dim)
			if err != nil {
				log.Fatalf("could not parse data in file %s: %s\n", path, err)
			}
			err = dim.validate()
			if err != nil {
				log.Fatalf("validation error in file %s: %s\n", path, err)
			}
			err = setDimension(ctx, doc, dim.Info.Id, raw)
			if err != nil {
				log.Fatalln(err)
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
		log.Debugln("Updating dimension " + dimensionID)
		err = dimension.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not update %s with %s: %s", "dimension", dimensionID, err)
		}
	} else {
		log.Debugln("Creating dimension " + dimensionID)
		_, err = doc.CreateDimensionRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not create %s with %s: %s", "dimension", dimensionID, err)
		}
	}
	return nil
}
