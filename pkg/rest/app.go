package rest

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/qlik-oss/corectl/pkg/log"
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

// TranslateAppNameToId translates a user provided app identifier into a app resourceId that can be used to open a websocket
// by checking against the items service if it is a resourceId, itemId or an app name
func (c *RestCaller) TranslateAppNameToId(userProvidedAppIdentifier string) string {

	var resourceIdBasedOnExplicitResourceId string
	var resourceIdBasedOnItemId string
	var resourceIdBasedOnName string

	// Check if we get a match in any of the three app identifier "formats" in parallel to reduce latency
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		resourceIdBasedOnExplicitResourceId = c.translateResourceIdToResourceId(userProvidedAppIdentifier)
		wg.Done()
	}()
	go func() {
		resourceIdBasedOnItemId = c.translateItemIdToResourceId(userProvidedAppIdentifier)
		wg.Done()
	}()
	go func() {
		resourceIdBasedOnName = c.translateAppNameToResourceId(userProvidedAppIdentifier)
		wg.Done()
	}()
	wg.Wait()

	if resourceIdBasedOnExplicitResourceId != "" {
		return resourceIdBasedOnExplicitResourceId
	}
	if resourceIdBasedOnItemId != "" {
		return resourceIdBasedOnItemId
	}
	if resourceIdBasedOnName != "" {
		return resourceIdBasedOnName
	}
	log.Fatalf("No such app found: %s", userProvidedAppIdentifier)
	return ""
}

func (c *RestCaller) translateItemIdToResourceId(potentialItemId string) string {
	var result RestDocListItem
	err := c.CallStd("GET", fmt.Sprintf("v1/items/%s", potentialItemId), "", nil, nil, &result)
	if err != nil {
		return ""
	}
	fmt.Println(result)
	return result.ResourceId
}

func (c *RestCaller) translateResourceIdToResourceId(potentialItemId string) string {
	var result ListAppResponse
	err := c.CallStd("GET", "v1/items", "", map[string]string{"resourceId": potentialItemId, "resourceType": "app"}, nil, &result)
	if err != nil {
		return ""
	}
	docList := result.Data
	if len(docList) > 1 {
		log.Fatalf("Too many apps matching the provided name: %s", potentialItemId)
		return ""
	}
	if len(docList) == 1 {
		return docList[0].ResourceId
	}
	return ""
}
func (c *RestCaller) translateAppNameToResourceId(potentialAppName string) string {
	var result ListAppResponse
	err := c.CallStd("GET", "v1/items", "", map[string]string{"sort": "-updatedAt", "limit": "30", "name": potentialAppName}, nil, &result)
	if err != nil {
		return ""
	}
	docList := result.Data
	if len(docList) > 29 { // There are possibly even more hits than we have in the response
		log.Fatalf("Too many apps matching the provided name: %s", potentialAppName)
		return ""
	}

	candidate := ""
	for _, x := range docList {
		if x.Name == potentialAppName {
			if candidate != "" {
				log.Fatalf("There are multiple apps matching the provided name: %s", potentialAppName)
				return ""
			}
			candidate = x.ResourceId
		}
	}
	return candidate
}

type ListAppResponse struct {
	Data []RestDocListItem `json:"data"`
}

type RestDocListItem struct {
	Name       string `json:"name"`
	ResourceId string `json:"resourceID"`
}
