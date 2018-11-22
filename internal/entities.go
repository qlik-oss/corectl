package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/qlik-oss/enigma-go"
	"gopkg.in/yaml.v2"
)

// EntitiesConfigFile defines a file that contains the different entities  yaml string arrays.
type EntitiesConfigFile struct {
	Measures   []string
	Dimensions []string
	Objects    []string
}
type genericEntity struct {
	Info *enigma.NxInfo `json:"qInfo"`
}

// ReadEntitiesFile reads the entity config file from the supplied path.
func ReadEntitiesFile(path string) EntitiesConfigFile {
	var config EntitiesConfigFile
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not find entities file:", path)
		os.Exit(1)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// SetupEntities reads all generic entities of the specified type from both the project file path and the config file path and updates
// the list of entities in the app.
func SetupEntities(ctx context.Context, doc *enigma.Doc, projectFile string, entitiesPathsOnCommandLine string, entityType string) {
	entitiesOnCommandLine, err := filepath.Glob(entitiesPathsOnCommandLine)
	if err != nil {
		FatalError(err)
	}
	for _, relativeEntityPath := range entitiesOnCommandLine {
		setupEntity(ctx, doc, relativeEntityPath, entityType)
	}
	if projectFile != "" {
		configFileContents := ReadEntitiesFile(projectFile)
		currentWorkingDir, _ := os.Getwd()
		defer os.Chdir(currentWorkingDir)
		os.Chdir(filepath.Dir(projectFile))

		var entities []string
		switch entityType {
		case "dimension":
			entities = configFileContents.Dimensions
		case "measure":
			entities = configFileContents.Measures
		case "object":
			entities = configFileContents.Objects
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
	var props genericEntity
	err = json.Unmarshal(entityFileContents, &props)
	validateEntity(props, entityPath, err)

	//I do not know how to to this nicer, with less duplication.
	switch entityType {
	case "dimension":
		dimension, err := doc.GetDimension(ctx, props.Info.Id)
		if err == nil && dimension.Handle != 0 {
			LogVerbose("Updating dimension " + props.Info.Id)
			err = dimension.SetPropertiesRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to update dimension "+props.Info.Id, err)
			}
		} else {
			LogVerbose("Creating dimension " + props.Info.Id)
			_, err = doc.CreateDimensionRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to create dimension "+props.Info.Id, err)
			}
		}
	case "measure":
		measure, err := doc.GetMeasure(ctx, props.Info.Id)
		if err == nil && measure.Handle != 0 {
			LogVerbose("Updating measure " + props.Info.Id)
			err = measure.SetPropertiesRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to update measure "+props.Info.Id, err)
			}
		} else {
			LogVerbose("Creating measure " + props.Info.Id)
			_, err = doc.CreateMeasureRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to create measure "+props.Info.Id, err)
			}
		}
	case "object":
		object, err := doc.GetObject(ctx, props.Info.Id)
		if err == nil && object.Handle != 0 {
			LogVerbose("Updating object " + props.Info.Id)
			err = object.SetPropertiesRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to update object "+props.Info.Id, err)
			}
		} else {
			LogVerbose("Creating object " + props.Info.Id)
			_, err = doc.CreateObjectRaw(ctx, entityFileContents)
			if err != nil {
				FatalError("Failed to create object "+props.Info.Id, err)
			}
		}
	}
}

func validateEntity(entity genericEntity, entityPath string, err error) {
	if err != nil {
		FatalError("Invalid json", err)
	}
	if entity.Info.Id == "" {
		FatalError("Missing qInfo qId attribute", entityPath)
	}
	if entity.Info.Type == "" {
		FatalError("Missing qInfo type attribute", entityPath)
	}
}
