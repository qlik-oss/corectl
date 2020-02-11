package rest

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
)

type RestCallerSettings interface {
	TlsConfig() *tls.Config
	Insecure() bool
	Headers() http.Header
	RestBaseUrl() *neturl.URL
	Engine() string
	EngineURL() *neturl.URL
	App() string
	AppId() string
}

type Caller struct {
	RestCallerSettings
}

// parseFunction is any function that reads bytes into a pointer and returns an error,
// json.Unmarshal is an example.
type parseFunction func([]byte, interface{}) error

// Call performs the specified the request and uses the passed parsing function 'read'
// to parse the response into the supplied result.
// If the response status code is not explicitly added to the map of accepted status codes
// an error will be returned.
func (c *Caller) Call(req *http.Request, result interface{}, statusCodes *map[int]bool, read parseFunction) error {
	client := http.DefaultClient
	certs := c.TlsConfig()
	if certs != nil {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: certs,
			},
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	ok := (*statusCodes)[response.StatusCode]
	if !ok {
		return fmt.Errorf("got status %d", response.StatusCode)
	}
	data, _ := ioutil.ReadAll(response.Body)
	err = read(data, result)
	if err != nil {
		return err
	}
	return nil
}

// Removes path and query escapes an app id.
func adaptAppID(appID string) string {
	split := strings.Split(appID, "/")
	adaptedID := split[len(split)-1]
	return neturl.QueryEscape(adaptedID)
}
