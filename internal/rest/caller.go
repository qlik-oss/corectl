package rest

import (
	"fmt"
	"io/ioutil"
	neturl "net/url"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// parseFunction is any function that reads bytes into a pointer and returns an error,
// json.Unmarshal is an example.
type parseFunction func([]byte, interface{}) error

// GetBaseURL returns the base URL for Rest API calls based on the value of 'engine'
func GetBaseURL() (*neturl.URL) {
	url := buildRestBaseURL(viper.GetString("engine"))
	return url
}

// Call performs the specified the request and uses the passed parsing function 'read'
// to parse the response into the supplied result.
// If the response status code is not explicitly added to the map of accepted status codes
// an error will be returned.
func Call(req *http.Request, result interface{}, statusCodes *map[int]bool, read parseFunction) error {
	response, err := http.DefaultClient.Do(req)
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

func buildRestBaseURL(engine string) (*neturl.URL) {
	u, _ := neturl.Parse(engine)
	switch (u.Scheme) {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
		default:
		url := "http://" + u.String()
		u, _ = neturl.Parse(url)
	}
	u.Path = ""
	return u
}

// Removes path and query escapes an app id.
func adaptAppID(appID string) string {
	split := strings.Split(appID, "/")
	adaptedID := split[len(split) - 1]
	return neturl.QueryEscape(adaptedID)
}
