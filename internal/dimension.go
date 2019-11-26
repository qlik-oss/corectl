package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// Dimension is a struct describing a generic dimension
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
	unsortedResult := make(map[string]*NamedItem)
	keys := make([]string, len(unsortedResult))
	for _, item := range layout.DimensionList.Items {
		parsedRawData := &ParsedEntityListData{}
		json.Unmarshal(item.Data, parsedRawData)
		unsortedResult[item.Info.Id] = &NamedItem{Title: parsedRawData.Title, ID: item.Info.Id}
		keys = append(keys, item.Info.Id)
	}

	//Loop over the keys that are sorted on qId and fetch the result for each object
	sort.Strings(keys)
	sortedResult := make([]NamedItem, len(keys))
	for i, key := range keys {
		sortedResult[i] = *unsortedResult[key]
	}
	return sortedResult
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
		ch := make(chan error)

		for _, raw := range rawEntities {
			go func(raw json.RawMessage) {
				var dim Dimension
				err := json.Unmarshal(raw, &dim)
				if err != nil {
					ch <- fmt.Errorf("could not parse data in file %s: %s", path, err)
					return
				}
				err = dim.validate()
				if err != nil {
					ch <- fmt.Errorf("validation error in file %s: %s", path, err)
					return
				}
				ch <- setDimension(ctx, doc, dim.Info.Id, raw)
			}(raw)
		}

		// Loop through the responses and see if there are any failures, if so exit with a fatal
		success := true
		for range rawEntities {
			err := <-ch
			if err != nil {
				log.Errorln(err)
				success = false
			}
		}

		if !success {
			log.Fatalln("One or more dimensions failed to be created or updated")
		}
	}
}

func setDimension(ctx context.Context, doc *enigma.Doc, dimensionID string, raw json.RawMessage) error {
	dimension, err := doc.GetDimension(ctx, dimensionID)
	if err != nil {
		return err
	}
	if dimension.Handle != 0 {
		log.Verboseln("Updating dimension " + dimensionID)
		err = dimension.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not update %s with %s: %s", "dimension", dimensionID, err)
		}
	} else {
		log.Verboseln("Creating dimension " + dimensionID)
		_, err = doc.CreateDimensionRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("could not create %s with %s: %s", "dimension", dimensionID, err)
		}
	}
	return nil
}
