package rest

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

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
// by checking if the provided identifier matches specific ID patterns. If it's not app ID,
// the items service will be used to get the resource ID for the provided identifier.
func (c *RestCaller) TranslateAppNameToId(userProvidedAppIdentifier string) string {

	var resourceID string

	if userProvidedAppIdentifier == "" {
		return ""
	}

	if isAppID(userProvidedAppIdentifier) {
		resourceID = userProvidedAppIdentifier
	} else if isItemID(userProvidedAppIdentifier) {
		resourceID = c.translateItemIdToResourceId(userProvidedAppIdentifier)
	} else {
		resourceID = c.translateAppNameToResourceId(userProvidedAppIdentifier)
	}

	if resourceID == "" {
		log.Fatalf("No such app found: %s", userProvidedAppIdentifier)
	}

	return resourceID
}

// isAppID returns true if the string 'x' is a lower-cased UUID, i.e. consisting of 5
// hyphenated hexadecimal numbers of different lenghts.
func isAppID(x string) bool {
	re := regexp.MustCompile("^[a-f0-9]+$")
	lengths := []int{8, 4, 4, 4, 12}
	parts := strings.Split(x, "-")
	for i, part := range parts {
		if len(part) != lengths[i] {
			return false
		} else if !re.MatchString(part) {
			return false
		}
	}
	return true
}

// isItemID returns true if the string 'x' is a 24-digit (lowercase) hexadecimal number.
func isItemID(x string) bool {
	re := regexp.MustCompile("^[a-f0-9]{24}$")
	return re.MatchString(x)
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
