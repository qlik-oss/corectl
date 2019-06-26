package internal

import (
	"net/url"
	"regexp"
	"strings"
)

// ParseEngineURL parses the engine parameter and returns an websocket URL if at all possible
func ParseEngineURL(engine string) (*url.URL) {
	if engine == "" {
		FatalError("engine URL not specified")
	}
	u, err := url.Parse(engine)
	if err != nil {
		// If 'engine' couldn't be parsed by net/url.Parse and contains colon
		// it is probably an IP address with a port then
		if strings.Contains(engine, ":") {
			// The passed url also contains at least one '/' so hopefully that is a path
			split := strings.Split(engine, "/")
			path := strings.Join(split[1:], "/")
			u = &url.URL{
				Host: split[0],
			}
			if len(path) > 0 {
				u.Path = "/" + path
			}
		}
	}
	switch u.Scheme {
	case "":
		if u.Host == "" {
			u, err = url.Parse("ws://" + u.String())
			if err != nil {
				FatalError("could not parse engine URL: ", engine)
			}
		} else {
			u.Scheme = "ws"
		}
	case "ws":
	case "wss":
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	default:
		// Strings of the form 'localhost:1234' are parsed as an URL with
		// Scheme = localhost, Opaque = 1234 which we want to turn into Host
		if u.Opaque != "" {
			u, err = url.Parse("ws://" + u.String())
			if err != nil {
				FatalError("could not parse engine URL: ", engine)
			}
		} else {
			FatalError("could not parse engine URL: ", engine)
		}
	}
	return u
}


func buildWebSocketURL(engine string, ttl string) string {
	u := ParseEngineURL(engine)
	// Only modify the URL path if there is no path set
	if u.Path == "" {
		u.Path = "/app/corectl/ttl/" + ttl
	}
	return u.String()
}

// TryParseAppFromURL parses an url for an app identifier
func TryParseAppFromURL(engineURL string) string {
	// Find any string in the path succeeding "/app/", and excluding anything after "/"
	re, _ := regexp.Compile("/app/([^/]+)")
	values := re.FindStringSubmatch(engineURL)
	if len(values) > 0 {
		appName := values[1]
		LogVerbose("Found app in engine url: " + appName)
		return appName
	}
	return ""
}
