package rest

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal/log"
	"net/http"
	"os"
	"strings"
)

// ImportApp imports a local app into the engine using the rest api
// To not have any dependency on internal, both appID and appName are returned.
func (c *RestCaller) ImportApp(appPath string) (appID, appName string, err error) {
	file, err := os.Open(appPath)
	if err != nil {
		err = fmt.Errorf("could not open file: %s", appPath)
		return
	}
	defer file.Close()

	header := make(http.Header)
	header.Add("Content-Type", "binary/octet-stream")
	req := &http.Request{
		Method: "POST",
		URL:    c.CreateUrl("v1/apps/import", nil),
		Header: header,
		Body:   file,
	}
	appInfo := &NxApp{}
	err = c.Call(req, appInfo, json.Unmarshal)
	if err != nil {
		err = fmt.Errorf("could not import app: %s", err.Error())
		return
	}
	appID = appInfo.Get("id")
	appName = appInfo.Get("name")
	return
}

func copy(originalMap http.Header) http.Header {
	newMap := make(http.Header)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

type NxApp struct {
	Attributes map[string]interface{} `json:"attributes"`
}

func (a NxApp) Get(attr string) string {
	attrVal := a.Attributes[attr]
	return fmt.Sprintf("%v", attrVal)
}

func (c *RestCaller) ListApps() ([]byte, error) {
	data, err := c.CallGet("v1/items", map[string]string{"sort": "-updatedAt", "limit": "30"})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *RestCaller) TranslateAppNameToId(name string) string {
	data, err := c.CallGet("v1/items", map[string]string{"sort": "-updatedAt", "limit": "30", "query": name})
	if err != nil {
		return ""
	}
	var result ListAppResponse
	json.Unmarshal(data, &result)
	docList := result.Data
	for _, x := range docList {
		if x.DocName == name {
			return x.DocId
		}
	}
	return ""
}

// PrintApps prints a list of apps and some meta to system out.
func PrintApps(data []byte, mode log.PrintMode) {
	if mode.JsonMode() {
		log.PrintAsJSON(data)
	} else {
		var result ListAppResponse
		json.Unmarshal(data, &result)
		docList := result.Data
		if mode.BashMode() {
			for _, app := range docList {
				PrintToBashComp(app.DocId)
			}
		} else if mode.QuietMode() {
			for _, app := range docList {
				PrintToBashComp(app.DocId)
			}
		} else {
			writer := tablewriter.NewWriter(os.Stdout)
			writer.SetAutoFormatHeaders(false)
			writer.SetHeader([]string{"Id", "Name"})
			for _, doc := range docList {
				writer.Append([]string{doc.DocId, doc.DocName})
			}
			writer.Render()
		}
	}
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

type ListAppResponse struct {
	Data []RestDocListItem `json:"data"`
}

type RestDocListItem struct {
	DocName string `json:"name"`
	DocId   string `json:"resourceID"`
}
