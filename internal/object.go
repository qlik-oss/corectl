package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qlik-oss/enigma-go"
)

type Object struct {
	Info       *enigma.NxInfo                  `json:"qInfo,omitempty"`
	Properties *enigma.GenericObjectProperties `json:"qProperty,omitempty"`
}

func (o Object) validate() error {
	if o.Info != nil {
		if o.Info.Id == "" {
			return errors.New("missing qInfo qId attribute")
		}
		if o.Info.Type == "" {
			return errors.New("missing qInfo qType attribute")
		}
	} else if o.Properties != nil {
		if o.Properties.Info == nil {
			return errors.New("missing qInfo attribute inside the qProperty")
		}
		if o.Properties.Info.Id == "" {
			return errors.New("missing qInfo qId attribute inside qProperty")
		}
		if o.Properties.Info.Type == "" {
			return errors.New("missing qInfo qType attribute inside qProperty")
		}
	} else {
		return errors.New("need to supply atleast one of qInfo or qProperty")
	}
	return nil
}

func ListObjects(ctx context.Context, doc *enigma.Doc) []NamedItemWithType {
	allInfos, _ := doc.GetAllInfos(ctx)
	unsortedResult := make(map[string]*NamedItemWithType)
	resultInOriginalOrder := []NamedItemWithType{}

	waitChannel := make(chan *NamedItemWithType)
	defer close(waitChannel)

	for _, item := range allInfos {
		go func(item *enigma.NxInfo) {
			object, _ := doc.GetObject(ctx, item.Id)
			if object != nil && object.Type != "" {
				rawProps, _ := object.GetPropertiesRaw(ctx)
				propsWithTitle := &PropsWithTitle{}
				json.Unmarshal(rawProps, propsWithTitle)
				waitChannel <- &NamedItemWithType{Title: propsWithTitle.Title, Id: item.Id, Type: item.Type}
			} else {
				waitChannel <- nil
			}
		}(item)
	}
	//Put all responses into a map by their Id
	for range allInfos {
		item := <-waitChannel
		if item != nil {
			unsortedResult[item.Id] = item
		}
	}
	//Loop over the original sort order, fetch the result items from the map and build the final result array
	for _, item := range allInfos {
		if unsortedResult[item.Id] != nil {
			resultInOriginalOrder = append(resultInOriginalOrder, *unsortedResult[item.Id])
		}
	}
	return resultInOriginalOrder
}

func SetObjects(ctx context.Context, doc *enigma.Doc, commandLineGlobPattern string) {
	paths, err := getEntityPaths(commandLineGlobPattern, "objects")
	if err != nil {
		FatalError("could not interpret glob pattern: ", err)
	}
	for _, path := range paths {
		rawEntities, err := parseEntityFile(path)
		if err != nil {
			FatalErrorf("could not parse file %s: %s", path, err)
		}
		for _, raw := range rawEntities {
			var object Object
			err = json.Unmarshal(raw, &object)
			if err != nil {
				FatalErrorf("could not parse data in file %s: %s", path, err)
			}
			err = object.validate()
			if err != nil {
				FatalErrorf("validation error in file %s: %s", path, err)
			}
			err = setObject(ctx, doc, object.Info, object.Properties, raw)
			if err != nil {
				FatalError(err)
			}
		}
	}
}

func setObject(ctx context.Context, doc *enigma.Doc, info *enigma.NxInfo, props *enigma.GenericObjectProperties, raw json.RawMessage) error {
	var objectID string
	isGenericObjectEntry := false
	if info != nil {
		objectID = info.Id
	} else {
		objectID = props.Info.Id
		isGenericObjectEntry = true
	}
	object, err := doc.GetObject(ctx, objectID)
	if err != nil {
		return err
	}
	if object.Handle != 0 {
		if isGenericObjectEntry {
			LogVerbose("Updating object " + objectID + " using SetFullPropertyTree")
			err = object.SetFullPropertyTreeRaw(ctx, raw)
		} else {
			LogVerbose("Updating object " + objectID + " using SetProperties")
			err = object.SetPropertiesRaw(ctx, raw)
		}
		if err != nil {
			return fmt.Errorf("failed to update %s %s: %s", "object", objectID, err)
		}
	} else {
		LogVerbose("Creating object " + objectID)
		if isGenericObjectEntry {
			var createdObject *enigma.GenericObject
			objectType := props.Info.Type
			createdObject, err = doc.CreateObject(ctx, &enigma.GenericObjectProperties{Info: &enigma.NxInfo{Id: objectID, Type: objectType}})
			LogVerbose("Setting object  " + objectID + " using SetFullPropertyTree")
			err = createdObject.SetFullPropertyTreeRaw(ctx, raw)
		} else {
			_, err = doc.CreateObjectRaw(ctx, raw)
		}
		if err != nil {
			return fmt.Errorf("failed to create %s %s: %s", "object", objectID, err)
		}
	}
	return nil
}
