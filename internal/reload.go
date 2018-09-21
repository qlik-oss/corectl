package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/qlik-oss/enigma-go"
)

// Reload reloads the app and prints the progress to system out. If true is supplied to skipTransientLogs
// the live ticking of table row counts is disabled (useful for testing).
func Reload(ctx context.Context, doc *enigma.Doc, global *enigma.Global, verboseLogging bool, skipTransientLogs bool) {

	var (
		reloadSuccessful bool
		err              error
	)

	// Only log reload progress if verbose logging is enabled
	if verboseLogging {
		reloadDone := make(chan struct{})
		ctxWithReservedRequestID, reservedRequestID := doc.WithReservedRequestID(ctx)
		go func() {
			for {
				select {
				case <-reloadDone:
					logProgress(ctx, global, reservedRequestID, skipTransientLogs)
					return
				default:
					time.Sleep(1000)
					// Get the progress using the request id we reserved for the reload
					logProgress(ctx, global, reservedRequestID, skipTransientLogs)
				}

			}
		}()
		reloadSuccessful, err = doc.DoReload(ctxWithReservedRequestID, 0, false, false)
		close(reloadDone)
	} else {
		reloadSuccessful, err = doc.DoReload(ctx, 0, false, false)
	}

	if err != nil {
		fmt.Println("Error when reloading app", err)
		panic(err)
	}
	if !reloadSuccessful {
		fmt.Println("DoReload was not successful!")
		os.Exit(1)
	}

	fmt.Println("Reload finished successfully")

}

func logProgress(ctx context.Context, global *enigma.Global, reservedRequestID int, skipTransientLogs bool) {
	progress, err := global.GetProgress(ctx, reservedRequestID)
	if err != nil {
		fmt.Println(err)
	} else {
		var text string
		if progress.TransientProgress != "" {
			if !skipTransientLogs {
				text = progress.TransientProgress
				fmt.Print("\033\r" + text)
			}
		} else if progress.PersistentProgress != "" {
			text = progress.PersistentProgress
			fmt.Print(text)
		}
	}
}

// Save calls DoSave on the app and prints "Done" if it succeeded or "Save failed" to system out.
func Save(ctx context.Context, doc *enigma.Doc, path string) {
	fmt.Print("Saving...")
	err := doc.DoSave(ctx, path)
	if err == nil {
		fmt.Println("Done")
	} else {
		fmt.Println("Save failed")
	}
	fmt.Println()

}
