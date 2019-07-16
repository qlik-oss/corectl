package internal

import (
	"context"

	"github.com/qlik-oss/enigma-go"
)

// ListAlternateState will return a lists of all alternate states in an app
func ListAlternateState(ctx context.Context, doc *enigma.Doc) []string {
	appLayout, _ := doc.GetAppLayout(ctx)
	return appLayout.StateNames
}

// AddAlternateState will add a named alternate state in the app
func AddAlternateState(ctx context.Context, doc *enigma.Doc, alternateStateName string) {
	err := doc.AddAlternateState(ctx, alternateStateName)
	if err != nil {
		FatalErrorf("could not add state %s: %s ", alternateStateName, err)
	}
}

// RemoveAlternateState will remove a named alternate state in the app
func RemoveAlternateState(ctx context.Context, doc *enigma.Doc, alternateStateName string) {
	states := ListAlternateState(ctx, doc)
	stateNameExists := Contains(states, alternateStateName)

	if !stateNameExists {
		FatalErrorf("no state with the name '%s' found in the app", alternateStateName)
	}

	err := doc.RemoveAlternateState(ctx, alternateStateName)
	if err != nil {
		FatalErrorf("could not remove state %s: %s ", alternateStateName, err)
	}
}
