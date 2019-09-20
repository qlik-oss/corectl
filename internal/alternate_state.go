package internal

import (
	"context"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// ListAlternateStates will return a list of all alternate states in an app.
func ListAlternateStates(ctx context.Context, doc *enigma.Doc) []string {
	appLayout, _ := doc.GetAppLayout(ctx)
	return appLayout.StateNames
}

// AddAlternateState will add a named alternate state in the app.
func AddAlternateState(ctx context.Context, doc *enigma.Doc, alternateStateName string) {
	err := doc.AddAlternateState(ctx, alternateStateName)
	if err != nil {
		log.Fatalf("could not add state %s: %s \n", alternateStateName, err)
	}
}

// RemoveAlternateState will remove a named alternate state in the app.
func RemoveAlternateState(ctx context.Context, doc *enigma.Doc, alternateStateName string) {
	states := ListAlternateStates(ctx, doc)
	var stateNameExists bool
	for _, state := range states {
		if state == alternateStateName {
			stateNameExists = true
			break
		}
	}

	if !stateNameExists {
		log.Fatalf("no alternate state with the name '%s' found in the app\n", alternateStateName)
	}

	err := doc.RemoveAlternateState(ctx, alternateStateName)
	if err != nil {
		log.Fatalf("could not remove state %s: %s \n", alternateStateName, err)
	}
}
