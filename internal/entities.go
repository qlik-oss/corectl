package internal

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

type genericEntity struct {
	Info     *enigma.NxInfo                  `json:"qInfo,omitempty"`
	Property *enigma.GenericObjectProperties `json:"qProperty,omitempty"`
}

// SetupEntities reads all generic entities of the specified type from both the project file path and the config file path and updates
// the list of entities in the app.
func SetupEntities(ctx context.Context, doc *enigma.Doc, entitiesPathsOnCommandLine string, entityType string) {
	entitiesOnCommandLine, err := filepath.Glob(entitiesPathsOnCommandLine)
	if err != nil {
		FatalError(err)
	}
	if len(entitiesOnCommandLine) > 0 {
		for _, relativeEntityPath := range entitiesOnCommandLine {
			setupEntity(ctx, doc, relativeEntityPath, entityType)
		}
	} else if ConfigDir != "" {
		currentWorkingDir, _ := os.Getwd()
		defer os.Chdir(currentWorkingDir)
		os.Chdir(ConfigDir)

		var entities []string
		switch entityType {
		case "dimension":
			entities = viper.GetStringSlice("dimensions")
		case "measure":
			entities = viper.GetStringSlice("measures")
		case "object":
			entities = viper.GetStringSlice("objects")
		default:
			FatalError("Unknown type: " + entityType)
		}

		for _, entityGlobPatternInConfigFile := range entities {
			entitiesInConfigLine, err := filepath.Glob(entityGlobPatternInConfigFile)
			if err != nil {
				FatalError(err)
			}
			for _, relativeEntityPath := range entitiesInConfigLine {
				setupEntity(ctx, doc, relativeEntityPath, entityType)
			}
		}
	}
}

func setupEntity(ctx context.Context, doc *enigma.Doc, entityPath string, entityType string) {
	entityFileContents, err := ioutil.ReadFile(entityPath)
	if err != nil {
		FatalError("Could not open "+entityType+" file", err)
	}
	var entity genericEntity
	err = json.Unmarshal(entityFileContents, &entity)
	validateEntity(entity, entityPath, err)
	//I do not know how to to this nicer, with less duplication.
	switch entityType {
	case "dimension":
		dimension, err := doc.GetDimension(ctx, entity.Info.Id)
		if err == nil && dimension.Handle != 0 {
			LogVerbose("Updating dimension " + entity.Info.Id)
			err = dimension.SetPropertiesRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to update dimension "+entity.Info.Id, err)
			}
		} else {
			LogVerbose("Creating dimension " + entity.Info.Id)
			_, err = doc.CreateDimensionRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to create dimension "+entity.Info.Id, err)
			}
		}
	case "measure":
		measure, err := doc.GetMeasure(ctx, entity.Info.Id)
		if err == nil && measure.Handle != 0 {
			LogVerbose("Updating measure " + entity.Info.Id)
			err = measure.SetPropertiesRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to update measure "+entity.Info.Id, err)
			}
		} else {
			LogVerbose("Creating measure " + entity.Info.Id)
			_, err = doc.CreateMeasureRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to create measure "+entity.Info.Id, err)
			}
		}
	case "object":
		var objectID string
		isGenericObjectEntry := false
		if entity.Info != nil {
			objectID = entity.Info.Id
		} else {
			objectID = entity.Property.Info.Id
			isGenericObjectEntry = true
		}
		object, err := doc.GetObject(ctx, objectID)
		if err == nil && object.Handle != 0 {
			if isGenericObjectEntry {
				LogVerbose("Updating object " + objectID + " using SetFullPropertyTree")
				err = object.SetFullPropertyTreeRaw(ctx, entityFileContents)
			} else {
				LogVerbose("Updating object " + objectID + " using SetProperties")
				err = object.SetPropertiesRaw(ctx, entityFileContents)
			}
			if err != nil {
				FatalError("Failed to update object "+objectID, err)
			}
		} else {
			LogVerbose("Creating object " + objectID)
			if isGenericObjectEntry {
				var createdObject *enigma.GenericObject
				createdObject, err = doc.CreateObject(ctx, &enigma.GenericObjectProperties{Info: &enigma.NxInfo{Id: objectID, Type: entity.Property.Info.Type}})
				LogVerbose("Setting object  " + objectID + " using SetFullPropertyTree")
				err = createdObject.SetFullPropertyTreeRaw(ctx, entityFileContents)
			} else {
				_, err = doc.CreateObjectRaw(ctx, entityFileContents)
			}
			if err != nil {
				FatalError("Failed to create object "+objectID, err)
			}
		}
	}
}

func validateEntity(entity genericEntity, entityPath string, err error) {
	if err != nil {
		FatalError("Invalid json", err)
	}
	if entity.Info == nil && entity.Property == nil {
		FatalError("Missing qInfo attribute or qProperty attribute", entityPath)
	}
	if entity.Info != nil && entity.Info.Id == "" {
		FatalError("Missing qInfo qId attribute", entityPath)
	}
	if entity.Info != nil && entity.Info.Type == "" {
		FatalError("Missing qInfo qType attribute", entityPath)
	}
	if entity.Property != nil && entity.Property.Info == nil {
		FatalError("Missing qInfo attribute inside the qProperty", entityPath)
	}
	if entity.Property != nil && entity.Property.Info.Id == "" {
		FatalError("Missing qInfo qId attribute inside qProperty", entityPath)
	}
	if entity.Property != nil && entity.Property.Info.Type == "" {
		FatalError("Missing qInfo qType attribute inside qProperty", entityPath)
	}
}
