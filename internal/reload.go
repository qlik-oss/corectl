package internal

import (
	"context"
	"time"

	"github.com/qlik-oss/enigma-go"
)

// Reload reloads the app and prints the progress to system out. If true is supplied to skipTransientLogs
// the live ticking of table row counts is disabled (useful for testing).
func Reload(ctx context.Context, doc *enigma.Doc, global *enigma.Global, silent bool, skipTransientLogs bool) {

	var (
		reloadSuccessful bool
		err              error
	)

	// Log progress unless silent flag was passed in to the reload command
	if !silent {
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
		//fetch the progress but do nothing, othwerwise we will get it for the next non silent call
		_, getProgressErr := global.GetProgress(ctx, 0)
		if getProgressErr != nil {
			Logger.Info(getProgressErr)
		}
	}

	if err != nil {
		Logger.Errorf("Error when reloading app: %s", err)
	}
	if !reloadSuccessful {
		Logger.Fatal("DoReload was not successful!")
	}

	Logger.Info("Reload finished successfully")

}

func logProgress(ctx context.Context, global *enigma.Global, reservedRequestID int, skipTransientLogs bool) {
	progress, err := global.GetProgress(ctx, reservedRequestID)
	if err != nil {
		Logger.Error(err)
	} else {
		var text string
		if progress.TransientProgress != "" {
			if !skipTransientLogs {
				text = progress.TransientProgress
				Logger.Info("\033\r" + text)
			}
		} else if progress.PersistentProgress != "" {
			text = progress.PersistentProgress
			Logger.Info(text)
		}
	}
}

// Save calls DoSave on the app and prints "Done" if it succeeded or "Save failed" to system out.
func Save(ctx context.Context, doc *enigma.Doc, path string) {
	Logger.Info("Saving...")
	err := doc.DoSave(ctx, path)
	if err == nil {
		Logger.Info("Done")
	} else {
		Logger.Error("Save failed")
	}
}
