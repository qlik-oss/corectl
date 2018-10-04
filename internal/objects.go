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
	if projectFile != "" {
		configFileContents := ReadObjectsFile(projectFile)
		for _, objectPathInConfigFile := range configFileContents.Objects {
			pathRelativeToProjectFile := RelativeToProject(projectFile, objectPathInConfigFile)
			setupObject(ctx, doc, pathRelativeToProjectFile)
		}
	}
	for _, relativeObjectPath := range filepath.SplitList(objectPathsOnCommandLine) {
		setupObject(ctx, doc, relativeObjectPath)
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
	object, _ := doc.GetObject(ctx, props.Info.Id)
	if object != nil {
		object.SetPropertiesRaw(ctx, objectFileContents)
	} else {
		doc.CreateObjectRaw(ctx, objectFileContents)
	}
}
