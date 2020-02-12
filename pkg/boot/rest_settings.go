package boot

import (
	neturl "net/url"
)

type RestSettings struct {
	*CommonSettings
}

func (c *RestSettings) RestBaseUrl() *neturl.URL {
	u := c.EngineURL()
	// CreateBaseURL returns the base URL for Rest API calls based on the value of 'engine'

	if u.Scheme == "ws" {
		u.Scheme = "http"
	} else if u.Scheme == "wss" {
		u.Scheme = "https"
	}
	return u
}
