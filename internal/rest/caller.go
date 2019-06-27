package rest

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
)

// parseFunction is any function that reads bytes into a pointer and returns an error,
// json.Unmarshal is an example.
type parseFunction func([]byte, interface{}) error

// CreateBaseURL returns the base URL for Rest API calls based on the value of 'engine'
func CreateBaseURL(u neturl.URL) *neturl.URL {
	if u.Scheme == "ws" {
		u.Scheme = "http"
	} else if u.Scheme == "wss" {
		u.Scheme = "https"
	}
	return &u
}

// Call performs the specified the request and uses the passed parsing function 'read'
// to parse the response into the supplied result.
// If the response status code is not explicitly added to the map of accepted status codes
// an error will be returned.
func Call(req *http.Request, certs *tls.Config, result interface{}, statusCodes *map[int]bool, read parseFunction) error {
	client := http.DefaultClient
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
