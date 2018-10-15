package printer

import (
	"encoding/json"
	"fmt"
	"github.com/qlik-oss/corectl/internal"
	"time"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/enigma-go"
)

// PrintApps prints a list of apps and some meta to system out.
func PrintApps(docList []*enigma.DocListEntry, printAsJSON bool) {
	if printAsJSON {
		buffer, err := json.Marshal(filterDocEntries(docList))
		if err != nil {
			internal.FatalError(err)
		}
		fmt.Println(prettyJSON(buffer))
	} else {
		docTable := tm.NewTable(0, 10, 3, ' ', 0)
		fmt.Fprintf(docTable, "Id\tName\tLast-Reloaded\tReadOnly\tTitle\n")
		for _, doc := range docList {
			fmt.Fprintf(docTable, "%s\t%s\t%s\t%t\t%s\n", doc.DocId, doc.DocName, doc.LastReloadTime, doc.ReadOnly, doc.Title)
		}
		fmt.Print(docTable)
	}
}

type filteredDocEntry struct {
	// Identifier of the app.
	DocID string `json:"id"`
	// Name of the app.
	DocName string `json:"name"`
	// Title of the app.
	Title string `json:"title"`
	// Last modified time stamp of the app.
	LastModifiedTime time.Time `json:"lastModifiedTime"`
	// Meta data related to the app.
	LastReloadTime string `json:"lastReloadTime"`
	// Size of remote app.
	FileSize enigma.Float64 `json:"fileSize"`
	// If set to true, the app is read-only.
	ReadOnly bool `json:"readOnly"`
}

func filterDocEntries(docList []*enigma.DocListEntry) []*filteredDocEntry {
	result := make([]*filteredDocEntry, len(docList))
	for i, doc := range docList {
		result[i] = &filteredDocEntry{
			DocName:          doc.DocName,
			LastModifiedTime: serialTimeToString(doc.FileTime),
			FileSize:         doc.FileSize,
			DocID:            doc.DocId,
			LastReloadTime:   doc.LastReloadTime,
			ReadOnly:         doc.ReadOnly,
			Title:            doc.Title,
		}
	}
	return result
}

func serialTimeToString(filetime enigma.Float64) time.Time {
	unix := (filetime - 25569) * 86400
	timestamp := time.Unix(int64(unix), 0).UTC()
	return timestamp
}
