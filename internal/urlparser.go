package internal

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/spf13/viper"
)

// GetEngineURL gets QIX engine URL from viper
func GetEngineURL() *url.URL {
	engine := viper.GetString("engine")
	if engine == "" {
		log.Fatalln("engine URL not specified")
	}
	u, err := parseEngineURL(engine)
	if err != nil {
		log.Fatalf("could not parse engine url '%s' got error: '%s'\n", engine, err)
	}
	return u
}

// parseEngineURL parses the engine parameter and returns an websocket URL if at all possible
// The following quirks/behavior of net/url.Parse might be nice to know
//
// 'localhost'                    => path=localhost
// 'localhost/app/foo'            => path=localhost/app/foo
// 'localhost:9076'               => scheme=localhost, opaque=9076
// 'localhost:9076/app/foo'       => scheme=localhost, opaque=9076/app/foo
// '127.0.0.1'                    => path=127.0.0.1
// '127.0.0.1:1234                => parse error
// 'ws://localhost:9076'          => scheme=ws, host=localhost:9076
// 'ws://example.com'             => scheme=ws, host=example.com
// 'ws://localhost:9076/app/foo'  => scheme=ws, host=localhost:9076, path=/app/foo
//
// References, if you need to see it with your own eyes:
//
// Documentation: https://golang.org/pkg/net/url/#URL
//        Source: https://github.com/golang/go/blob/master/src/net/url/url.go
//
func parseEngineURL(engine string) (*url.URL, error) {
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
				return nil, err
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
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return u, nil
}

func buildWebSocketURL(ttl string) string {
	u := GetEngineURL()
	// Only modify the URL path if there is no path set
	if u.Path == "" || u.Path == "/" {
		u.Path = "/app/engineData/ttl/" + ttl
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
		log.Verboseln("Found app in engine url: " + appName)
		return appName
	}
	return ""
}
