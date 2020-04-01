package boot

import (
	"crypto/tls"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/pkg/log"
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
	headers := c.GetHeaders()

	//TODO headers.Set("User-Agent", fmt.Sprintf("corectl/%s (%s)", version, runtime.GOOS))
	headers.Set("User-Agent", "corectl")
	return headers
}

func (c *CommonSettings) Server() string {
	return c.GetString("server")
}

// serverURL returns the parsed engine url
func (c *CommonSettings) serverURL() *neturl.URL {
	server := c.Server()
	if server == "" {
		log.Fatalln("server URL not specified")
	}
	u, err := parseURL(server, "ws")
	if err != nil {
		log.Fatalf("could not parse server url '%s' got error: '%s'\n", server, err)
	}
	return u
}

////////////////////////////////////////
// Functions for connecting over rest //
////////////////////////////////////////

func (c *CommonSettings) IsSenseForKubernetes() bool {
	if c.serverURL().Scheme == "http" || c.serverURL().Scheme == "https" {
		return true
	} else {
		return false
	}
}

func (c *CommonSettings) RestBaseUrl() *neturl.URL {
	u, _ := neturl.Parse(c.serverURL().String()) //Clone it since we are going to modify it
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
	engineUrl := c.serverURL()
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
	serverURL := c.serverURL()
	return serverURL.Host
}

func (c *CommonSettings) App() string {
	appName := c.GetString("app")
	if appName == "" {
		appName = tryParseAppFromURL(c.GetString("server"))
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
