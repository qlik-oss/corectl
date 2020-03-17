package rest

import (
	"fmt"
	"net/http"
	"os"
)

// ImportApp imports a local app into the engine using the rest api
// To not have any dependency on internal, both appID and appName are returned.
func (c *RestCaller) ImportApp(appPath string) (appID, appName string, err error) {
	if c.IsSenseForKubernetes() {
		log.Fatalln("Not implemented for Sense yet")
	}

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
	err = c.CallReq(req, appInfo)
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
	var result []byte
	err := c.CallStd("GET", "v1/items", "", map[string]string{"sort": "-updatedAt", "limit": "30"}, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RestCaller) TranslateAppNameToId(name string) string {
	var result ListAppResponse
	err := c.CallStd("GET", "v1/items", "", map[string]string{"sort": "-updatedAt", "limit": "30", "query": name}, nil, &result)
	if err != nil {
		return ""
	}
	docList := result.Data
	for _, x := range docList {
		if x.DocName == name {
			return x.DocId
		}
	}
	return ""
}

type ListAppResponse struct {
	Data []RestDocListItem `json:"data"`
}

type RestDocListItem struct {
	DocName string `json:"name"`
	DocId   string `json:"resourceID"`
}
