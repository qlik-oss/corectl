package boot

import (
	"crypto/tls"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"net/http"
	neturl "net/url"
	"strings"
)

func NewCommonSettings(cfg *dynconf.DynSettings) *CommonSettings {
	commonSettings := &CommonSettings{DynSettings: cfg}

	return commonSettings
}

type PrintMode struct {
	cfg *CommonSettings
}

func (o *PrintMode) QuietMode() bool {
	return o.cfg.GetBoolAllowNoFlag("quiet")
}
func (o *PrintMode) BashMode() bool {
	return o.cfg.GetBool("bash")

}
func (o *PrintMode) JsonMode() bool {
	return o.cfg.GetBool("json")
}
func (o *PrintMode) TrafficMode() bool {
	return o.cfg.GetBool("traffic")
}
func (o *PrintMode) VerboseMode() bool {
	return o.cfg.GetBool("verbose")
}

type CommonSettings struct {
	*dynconf.DynSettings
	ExplicitlyTranslatedAppId string
}

func (c *CommonSettings) PrintMode() log.PrintMode {
	return &PrintMode{cfg: c}
}

func (c *CommonSettings) TlsConfig() *tls.Config {
	tlsClientConfig := c.GetTLSConfigFromPath("certificates")
	if c.GetBool("insecure") {
		tlsClientConfig.InsecureSkipVerify = true
	}
	return tlsClientConfig
}

func (c *CommonSettings) Insecure() bool {
	return c.GetBool("insecure")
}

func (c *CommonSettings) Headers() http.Header {
	headersMap := c.GetStringMap("headers")

	var headers http.Header
	headers = make(http.Header, 1)
	for key, value := range headersMap {
		headers.Set(key, value)
	}
	//TODO headers.Set("User-Agent", fmt.Sprintf("corectl/%s (%s)", version, runtime.GOOS))
	headers.Set("User-Agent", "corectl")
	return headers
}

func (c *CommonSettings) Engine() string {
	return c.GetString("engine")
}

// engineURL returns the parsed engine url
func (c *CommonSettings) engineURL() *neturl.URL {
	engine := c.Engine()
	if engine == "" {
		log.Fatalln("engine URL not specified")
	}
	u, err := parseURL(engine, "ws")
	if err != nil {
		log.Fatalf("could not parse engine url '%s' got error: '%s'\n", engine, err)
	}
	return u
}

////////////////////////////////////////
// Functions for connecting over rest //
////////////////////////////////////////

func (c *CommonSettings) IsSenseForKubernetes() bool {
	if c.engineURL().Scheme == "http" || c.engineURL().Scheme == "https" {
		return true
	} else {
		return false
	}
}

func (c *CommonSettings) RestBaseUrl() *neturl.URL {
	u, _ := neturl.Parse(c.engineURL().String()) //Clone it since we are going to modify it
	// CreateBaseURL returns the base URL for Rest API calls based on the value of 'engine'
	if u.Scheme == "ws" {
		u.Scheme = "http"
	} else if u.Scheme == "wss" {
		u.Scheme = "https"
	}
	if c.IsSenseForKubernetes() {
		u.Path = "/api/"
	}
	return u
}

func (c *CommonSettings) RestAdaptedAppId() string {
	appId := c.AppId()
	split := strings.Split(appId, "/")
	adaptedID := split[len(split)-1]
	return neturl.QueryEscape(adaptedID)
}

/////////////////////////////////////////////
// Functions for connecting over websocket //
/////////////////////////////////////////////

func (c *CommonSettings) WebSocketEngineURL() string {
	engineUrl := c.engineURL()
	if c.IsSenseForKubernetes() {
		if engineUrl.Scheme == "https" {
			return "wss://" + engineUrl.Host + "/app/" + c.AppId()
		} else {
			return "ws://" + engineUrl.Host + "/app/" + c.AppId()
		}
	} else {
		// Only modify the URL path if there is no path set
		if engineUrl.Path == "" || engineUrl.Path == "/" {
			return engineUrl.Scheme + "://" + engineUrl.Host + "/app/engineData/ttl/" + c.Ttl()
		} else {
			return engineUrl.String()
		}
	}
}

func (c *CommonSettings) AppIdMappingNamespace() string {
	engineURL := c.engineURL()
	return engineURL.Host
}

func (c *CommonSettings) App() string {
	appName := c.GetString("app")
	if appName == "" {
		appName = tryParseAppFromURL(c.GetString("engine"))
	}
	return appName
}

func (c *CommonSettings) AppId() string {
	if c.ExplicitlyTranslatedAppId != "" {
		return c.ExplicitlyTranslatedAppId
	}
	appName := c.App()
	appId, _ := ApplyNameToIDTransformation(c.AppIdMappingNamespace(), appName)
	return appId
}

func (c *CommonSettings) Ttl() string {
	return c.GetString("ttl")
}

func (c *CommonSettings) NoData() bool {
	return c.GetBool("no-data")
}
func (c *CommonSettings) NoSave() bool {
	return c.GetBoolAllowNoFlag("no-save")
}
