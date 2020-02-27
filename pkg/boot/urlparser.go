package boot

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/qlik-oss/corectl/internal/log"
)

// tryParseAppFromURL parses an url for an app identifier
func tryParseAppFromURL(engineURL string) string {
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

func parseURL(engine string, defaultProtocol string) (*url.URL, error) {
	if !hasProtocol(engine) {
		engine = defaultProtocol + "://" + engine
	}
	return url.Parse(engine)
}

func hasProtocol(engine string) bool {
	lc := strings.ToLower(engine)
	switch {
	case strings.HasPrefix(lc, "ws://"):
		fallthrough
	case strings.HasPrefix(lc, "wss://"):
		fallthrough
	case strings.HasPrefix(lc, "http://"):
		fallthrough
	case strings.HasPrefix(lc, "https://"):
		return true
	default:
		return false
	}
}
