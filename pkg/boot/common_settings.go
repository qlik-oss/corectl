package boot

import (
	"crypto/tls"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"net/http"
	neturl "net/url"
)

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
}

func (n *CommonSettings) PrintMode() *PrintMode {
	return &PrintMode{cfg: n}
}

func (n *CommonSettings) TlsConfig() *tls.Config {
	tlsClientConfig := n.GetTLSConfigFromPath("certificates")
	if n.GetBool("insecure") {
		tlsClientConfig.InsecureSkipVerify = true
	}
	return tlsClientConfig
}

func (n *CommonSettings) Insecure() bool {
	return n.GetBool("insecure")
}

func (n *CommonSettings) Headers() http.Header {
	headersMap := n.GetStringMap("headers")

	var headers http.Header
	headers = make(http.Header, 1)
	for key, value := range headersMap {
		headers.Set(key, value)
	}
	//TODO headers.Set("User-Agent", fmt.Sprintf("corectl/%s (%s)", version, runtime.GOOS))
	return headers
}

func (n *CommonSettings) Engine() string {
	return n.GetString("engine")
}

// GetEngineURL gets QIX engine URL from viper
func (n *CommonSettings) EngineURL() *neturl.URL {
	engine := n.Engine()
	if engine == "" {
		log.Fatalln("engine URL not specified")
	}
	u, err := parseEngineURL(engine)
	if err != nil {
		log.Fatalf("could not parse engine url '%s' got error: '%s'\n", engine, err)
	}
	return u
}

func (n *CommonSettings) EngineHost() string {
	engineURL := n.EngineURL()
	return engineURL.Host
}

func (n *CommonSettings) App() string {
	specifiedApp := n.GetString("app")
	if specifiedApp != "" {
		return specifiedApp
	}
	return tryParseAppFromURL(n.EngineURL().Path)
}

func (n *CommonSettings) AppId() string {
	appName := n.App()
	if appName == "" {
		appName = tryParseAppFromURL(n.Engine())
	}
	appId, _ := ApplyNameToIDTransformation(n.EngineHost(), appName)
	return appId
}
