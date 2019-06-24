package rest

import (
	"encoding/json"
	"fmt"
	"os"
	neturl "net/url"
	"net/http"
)

// ImportApp imports a local app into the engine using the rest api
// To not have any dependency on internal, both appID and appName are returned.
func ImportApp(appPath string, headers http.Header) (appID, appName string, err error) {
	url := GetBaseURL()
	url.Path = "/v1/apps/import"
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
		URL: url,
		Header: headers,
		Body: file,
	}
	appInfo := &RestNxApp{}
	statusCodes := &map[int]bool{
		200: true,
	}
	err = Call(req, appInfo, statusCodes, json.Unmarshal)
	if err != nil {
		err = fmt.Errorf("could not import app: %s", err.Error())
		return
	}
	appID = appInfo.Get("id")
	appName = appInfo.Get("name")
	return
}

type RestNxApp struct {
	Attributes map[string]interface{} `json:"attributes"`
}

func (a RestNxApp) Get(attr string) string {
	attrVal := a.Attributes[attr]
	return fmt.Sprintf("%v", attrVal)
}
