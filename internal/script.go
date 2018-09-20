package internal

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/qlik-oss/enigma-go"
)

// SetScript loads the script file and sets it in the app.
func SetScript(ctx context.Context, doc *enigma.Doc, scriptFilePath string) {
	loadScript, err := ioutil.ReadFile(scriptFilePath)
	if err != nil {
		log.Fatalln("Could not find load script: %s", scriptFilePath)
	}
	err = doc.SetScript(ctx, string(loadScript))
}
