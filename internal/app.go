package internal

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// SetAppProperties loads the app properties from file and sets it in the app.
func SetAppProperties(ctx context.Context, doc *enigma.Doc, appProppertiesFilePath string) {
	content, err := ioutil.ReadFile(appProppertiesFilePath)
	if err != nil {
		log.Fatalf("could not find app-properties file: %s\n", appProppertiesFilePath)
	}
	var appProperties *enigma.NxAppProperties
	err = json.Unmarshal(content, &appProperties)
	if err != nil {
		log.Fatalf("could not parse app-properties in file %s: %s\n", appProppertiesFilePath, err)
	}

	err = doc.SetAppProperties(ctx, appProperties)

	if err != nil {
		log.Fatalln("failed to set app-properties: ", err)
	}
}
