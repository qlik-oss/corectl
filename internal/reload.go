package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

var transientLogged bool

// Reload reloads the app and prints the progress to system out.
func Reload(ctx context.Context, doc *enigma.Doc, global *enigma.Global, silent bool, limit int) {

	var (
		reloadSuccessful  bool
		err               error
		skipTransientLogs bool
	)

	// Log progress unless silent flag was passed in to the reload command
	if !silent {

		// If not running in a terminal we should skip transient progress logging
		skipTransientLogs = !terminal.IsTerminal(int(os.Stdout.Fd()))

		reloadDone := make(chan struct{})
		loggingDone := make(chan struct{})
		ctxWithReservedRequestID, reservedRequestID := doc.WithReservedRequestID(ctx)
		go func() {
			for {
				select {
				case <-reloadDone:
					logProgress(ctx, global, reservedRequestID, skipTransientLogs)
					close(loggingDone)
					return
				case <-time.After(time.Second):
					// Get the progress using the request id we reserved for the reload
					logProgress(ctx, global, reservedRequestID, skipTransientLogs)
				}
			}
		}()
		if limit == 0 {
			reloadSuccessful, err = doc.DoReload(ctxWithReservedRequestID, 0, false, false)
		} else {
			doc.SetFetchLimit(ctx, limit)
			reloadSuccessful, err = doc.DoReload(ctxWithReservedRequestID, 0, false, true)
		}

		close(reloadDone)
		<-loggingDone
	} else {
		reloadSuccessful, err = doc.DoReload(ctx, 0, false, false)
		//fetch the progress but do nothing, othwerwise we will get it for the next non silent call
		_, getProgressErr := global.GetProgress(ctx, 0)
		if getProgressErr != nil {
			fmt.Println(getProgressErr)
		}
	}

	if err != nil {
		FatalError("could not reload app: ", err)
	}
	if !reloadSuccessful {
		FatalError("reload was not successful")
	}

	fmt.Println("Reload finished successfully")

}

func logProgress(ctx context.Context, global *enigma.Global, reservedRequestID int, skipTransientLogs bool) {
	InteractDef := enigma.InteractDef{
		Type:      "IT_SCRIPTLINE",
		Title:     "",
		Msg:       "",
		Buttons:   0,
		Line:      "",
		OldLineNr: 0,
		NewLineNr: 0,
		Path:      "",
		Hidden:    false,
		Result:    0,
		Input:     "",
	}

	progress, err := global.GetProgress(ctx, reservedRequestID)
	if err != nil {
		fmt.Println(err)
	} else {
		// If debug is set just "OK" interaction
		if progress.UserInteractionWanted {
			global.InteractDone(ctx, reservedRequestID, &InteractDef)
		}

		var text string
		if progress.TransientProgress != "" {
			if !skipTransientLogs {
				text = progress.TransientProgress
				fmt.Print("\r" + text)
				transientLogged = true
			}
		} else if progress.PersistentProgress != "" {
			text = progress.PersistentProgress
			// If a transient progress was logged we should update that progress with the persistent one
			if transientLogged {
				fmt.Print("\r" + text)
				transientLogged = false
			} else {
				fmt.Print(text)
			}
		}
	}
}

// Save calls DoSave on the app and prints "Done" if it succeeded or "Save failed" to system out.
func Save(ctx context.Context, doc *enigma.Doc) {
	noData := viper.GetBool("no-data")
	var err error

	// If app is opened without data we should only save the objects
	if noData {
		fmt.Print("Saving objects in app... ")
		err = doc.SaveObjects(ctx)
	} else {
		fmt.Print("Saving app... ")
		err = doc.DoSave(ctx, "")
	}
	if err == nil {
		fmt.Println("Done")
	} else {
		fmt.Println("Save failed")
	}
	fmt.Println()
}
