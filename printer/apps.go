package printer

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/pkg/rest"
	"github.com/qlik-oss/enigma-go"
)

// PrintApps prints a list of apps and some meta to system out.
func PrintApps(docList []*enigma.DocListEntry, mode log.PrintMode) {
	if mode.JsonMode() {
		PrintAsJSON(filterDocEntries(docList))
	} else if mode.BashMode() {
		for _, app := range docList {
			PrintToBashComp(app.DocName)
		}
	} else if mode.QuietMode() {
		for _, app := range docList {
			PrintToBashComp(app.DocId)
		}
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Id", "Name", "Last-Reloaded", "ReadOnly", "Title"})
		for _, doc := range docList {
			writer.Append([]string{doc.DocId, doc.DocName, doc.LastReloadTime, strconv.FormatBool(doc.ReadOnly), doc.Title})
		}
		writer.Render()
	}
}

// PrintAppsRest prints a list of apps (from a REST call) and some meta to system out.
func PrintAppsRest(data []byte, mode log.PrintMode) {
	if mode.JsonMode() {
		PrintAsJSON(data)
	} else {
		var result rest.ListAppResponse
		json.Unmarshal(data, &result)
		docList := result.Data
		if mode.BashMode() {
			for _, app := range docList {
				PrintToBashComp(app.ResourceId)
			}
		} else if mode.QuietMode() {
			for _, app := range docList {
				PrintToBashComp(app.ResourceId)
			}
		} else {
			writer := tablewriter.NewWriter(os.Stdout)
			writer.SetAutoFormatHeaders(false)
			writer.SetHeader([]string{"Id", "Name"})
			for _, doc := range docList {
				writer.Append([]string{doc.ResourceId, doc.Name})
			}
			writer.Render()
		}
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
