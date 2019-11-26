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

// Measure is a struct describing a generic measure
type Measure struct {
	Info *enigma.NxInfo `json:"qInfo,omitempty"`
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

// ListMeasures fetches all measures and returns them in an array
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

	unsortedResult := make(map[string]*NamedItem)
	keys := make([]string, len(unsortedResult))
	for _, item := range layout.MeasureList.Items {
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

// SetMeasures creates or updates all measures on given glob patterns
func SetMeasures(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "measures")
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
				var measure Measure
				err := json.Unmarshal(raw, &measure)
				if err != nil {
					ch <- fmt.Errorf("could not parse data in file %s: %s", path, err)
					return
				}
				err = measure.validate()
				if err != nil {
					ch <- fmt.Errorf("validation error in file %s: %s", path, err)
					return
				}
				ch <- setMeasure(ctx, doc, measure.Info.Id, raw)
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
			log.Fatalln("One or more measures failed to be created or updated")
		}
	}
}

func setMeasure(ctx context.Context, doc *enigma.Doc, measureID string, raw json.RawMessage) error {
	measure, err := doc.GetMeasure(ctx, measureID)
	if err != nil {
		return err
	}
	if measure.Handle != 0 {
		log.Verboseln("Updating measure " + measureID)
		err = measure.SetPropertiesRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("failed to update %s with %s: %s", "measure", measureID, err)
		}
	} else {
		log.Verboseln("Creating measure " + measureID)
		_, err = doc.CreateMeasureRaw(ctx, raw)
		if err != nil {
			return fmt.Errorf("failed to create %s with %s: %s", "measure", measureID, err)
		}
	}
	return nil
}
