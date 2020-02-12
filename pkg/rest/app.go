package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
)

// ImportApp imports a local app into the engine using the rest api
// To not have any dependency on internal, both appID and appName are returned.
func (c *Caller) ImportApp(appPath string) (appID, appName string, err error) {
	url := c.RestBaseUrl()
	url.Path = "/v1/apps/import"
	headers := c.Headers()
	headers.Add("Content-Type", "binary/octet-stream")
	values := neturl.Values{}
	url.RawQuery = values.Encode()
	file, err := os.Open(appPath)
	if err != nil {
		err = fmt.Errorf("could not open file: %s", appPath)
		return
	}
	defer file.Close()
	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: headers,
		Body:   file,
	}
	appInfo := &NxApp{}
	statusCodes := &map[int]bool{
		200: true,
	}
	err = c.Call(req, appInfo, statusCodes, json.Unmarshal)
	if err != nil {
		err = fmt.Errorf("could not import app: %s", err.Error())
		return
	}
	appID = appInfo.Get("id")
	appName = appInfo.Get("name")
	return
}

type NxApp struct {
	Attributes map[string]interface{} `json:"attributes"`
}

func (a NxApp) Get(attr string) string {
	attrVal := a.Attributes[attr]
	return fmt.Sprintf("%v", attrVal)
}
