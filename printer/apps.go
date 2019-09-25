package printer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// PrintApps prints a list of apps and some meta to system out.
func PrintApps(docList []*enigma.DocListEntry, printAsBash bool) {
	switch mode {
	case jsonMode:
		log.PrintAsJSON(filterDocEntries(docList))
	case bashMode:
		for _, app := range docList {
			PrintToBashComp(app.DocName)
		}
	case quietMode:
		for _, app := range docList {
			PrintToBashComp(app.DocId)
		}
	default:
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Id", "Name", "Last-Reloaded", "ReadOnly", "Title"})
		for _, doc := range docList {
			writer.Append([]string{doc.DocId, doc.DocName, doc.LastReloadTime, strconv.FormatBool(doc.ReadOnly), doc.Title})
		}
		writer.Render()
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

// PrintToBashComp handles strings that should be included as options when using auto completion
func PrintToBashComp(str string) {
	if strings.Contains(str, " ") {
		// If string includes whitespaces we need to add quotes
		fmt.Printf("%q\n", str)
	} else {
		fmt.Println(str)
	}
}
