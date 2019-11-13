package internal

import (
	"context"
	"io/ioutil"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// SetScript loads the script file and sets it in the app.
func SetScript(ctx context.Context, doc *enigma.Doc, scriptFilePath string) {
	loadScript, err := ioutil.ReadFile(scriptFilePath)
	if err != nil {
		log.Fatalf("could not find loadscript: %s\n", scriptFilePath)
	}

	err = doc.SetScript(ctx, string(loadScript))

	if err != nil {
		log.Fatalln("failed to set script: ", err)
	}
}
