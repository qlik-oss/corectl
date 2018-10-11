package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qlik-oss/enigma-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ObjectsConfigFile defines a file that contains a list of objects as a yaml string array.
type ObjectsConfigFile struct {
	Objects []string
}

// ReadObjectsFile reads the object config file from the supplied path.
func ReadObjectsFile(path string) ObjectsConfigFile {
	var config ObjectsConfigFile
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not find objects file:", path)
		os.Exit(1)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// SetupObjects reads all objects from both the project file path and the config file path and updates
// the list of objects in the app.
func SetupObjects(ctx context.Context, doc *enigma.Doc, projectFile string, objectPathsOnCommandLine string) {

	objectsOnCommandLine, err := filepath.Glob(objectPathsOnCommandLine)
	if err != nil {
		FatalError(err)
	}
	for _, relativeObjectPath := range objectsOnCommandLine {
		setupObject(ctx, doc, relativeObjectPath)
	}
	if projectFile != "" {
		configFileContents := ReadObjectsFile(projectFile)
		currentWorkingDir, _ := os.Getwd()
		defer os.Chdir(currentWorkingDir)
		os.Chdir(filepath.Dir(projectFile))

		for _, objectGlobPatternInConfigFile := range configFileContents.Objects {
			objectsInConfigLine, err := filepath.Glob(objectGlobPatternInConfigFile)
			if err != nil {
				FatalError(err)
			}
			for _, relativeObjectPath := range objectsInConfigLine {
				setupObject(ctx, doc, relativeObjectPath)
			}
		}
	}
}

func setupObject(ctx context.Context, doc *enigma.Doc, objectPath string) {

	objectFileContents, err := ioutil.ReadFile(objectPath)
	if err != nil {
		FatalError("Could not open object file", err)
	}
	var props enigma.GenericObjectProperties
	err = json.Unmarshal(objectFileContents, &props)
	if err != nil {
		FatalError("Invalid json", err)
	}
	if props.Info.Id == "" {
		FatalError("Missing qInfo qId attribute", objectPath)
	}
	if props.Info.Type == "" {
		FatalError("Missing qInfo qId attribute", objectPath)
	}
	object, err := doc.GetObject(ctx, props.Info.Id)
	if err == nil && object.Handle != 0 {
		LogVerbose("Updating object " + props.Info.Id)
		err = object.SetPropertiesRaw(ctx, objectFileContents)
		if err != nil {
			FatalError("Failed to update object", err)
		}
	} else {
		LogVerbose("Creating object " + props.Info.Id)
		_, err = doc.CreateObjectRaw(ctx, objectFileContents)
		if err != nil {
			FatalError("Failed to create object", err)
		}
	}
}
