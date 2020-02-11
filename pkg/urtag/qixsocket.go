package urtag

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qlik-oss/corectl/pkg/huggorm"
	"github.com/qlik-oss/corectl/pkg/rest"
	"github.com/spf13/cobra"

	//"github.com/qlik-oss/corectl/printer"
	"net/http"
	"os"
	"os/user"
	//"runtime"
	"strconv"
	"strings"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

type EngineWebSocketSettings interface {
	TlsConfig() *tls.Config
	Insecure() bool
	Headers() http.Header

	Engine() string
	App() string
	AppId() string
	Ttl() string
	NoData() bool
	EngineHost() string
	WebSocketEngineURL() string
	WithoutApp() bool
	CreateAppIfMissing() bool
}

func NewCommunicator(ccmd *cobra.Command) *Communicator {
	cfg := huggorm.ReadSettings(ccmd)
	commonSettings := &CommonSettings{DynSettings: cfg}
	log.Init(commonSettings.PrintMode())
	return &Communicator{CommonSettings: commonSettings}
}

type Communicator struct {
	*CommonSettings
}

func (c *Communicator) RestCaller() *rest.Caller {
	restSettings := &RestSettings{CommonSettings: c.CommonSettings}
	return &rest.Caller{RestCallerSettings: restSettings}
}

func (c *Communicator) OpenAppSocket(createAppIfMissing bool) (context.Context, *enigma.Global, *enigma.Doc, *SocketSettings) {
	socketSettings := &SocketSettings{CommonSettings: c.CommonSettings, withoutApp: false, createAppifMissing: createAppIfMissing}
	global := GetGlobal(context.Background(), socketSettings)
	app := GetApp(context.Background(), global, socketSettings)
	return context.Background(), global, app, socketSettings
}
func (c *Communicator) OpenGlobal() (context.Context, *enigma.Global, *SocketSettings) {
	socketSettings := &SocketSettings{CommonSettings: c.CommonSettings, withoutApp: false, createAppifMissing: true}
	global := GetGlobal(context.Background(), socketSettings)
	return context.Background(), global, socketSettings
}

//AppExists returns wether or not an app exists along with any eventual error from engine.
func (c *Communicator) AppExists() (bool, error) {
	socketSettings := &SocketSettings{CommonSettings: c.CommonSettings, withoutApp: false, createAppifMissing: false}
	global := GetGlobal(context.Background(), socketSettings)
	_, err := global.GetAppEntry(context.Background(), socketSettings.AppId())
	if err != nil {
		return false, fmt.Errorf("could not find any app by ID '%s': %s", socketSettings.AppId(), err)
	}
	return true, nil
}

//DeleteApp removes the specified app from the engine.
func (c *Communicator) DeleteApp(appName string) {
	socketSettings := &SocketSettings{CommonSettings: c.CommonSettings, withoutApp: false, createAppifMissing: false}
	socketSettings.OverrideSetting("app", appName)
	global := GetGlobal(context.Background(), socketSettings)
	appID := socketSettings.AppId()
	succ, err := global.DeleteApp(context.Background(), socketSettings.AppId())
	if err != nil {
		log.Fatalf("could not delete app with name '%s' and ID '%s': %s\n", appName, appID, err)
	} else if !succ {
		log.Fatalf("could not delete app with name '%s' and ID '%s'\n", appName, appID)
	}
	SetAppIDToKnownApps(socketSettings.EngineHost(), appName, appID, true)
}

func (c *Communicator) RestClient() {

}

// State contains all needed info about the current app including a go context to use when communicating with the engine.
type State struct {
	doc     *enigma.Doc
	Ctx     context.Context
	Global  *enigma.Global
	AppName string
	AppID   string
	Verbose bool
}

func GetApp(ctx context.Context, global *enigma.Global, settings EngineWebSocketSettings) *enigma.Doc {
	engineHost := settings.EngineHost()
	appName := settings.App()

	if appName == "" {
		log.Fatalln("no app specified")
	}

	appID := settings.AppId()
	noData := settings.NoData()
	createAppIfMissing := settings.CreateAppIfMissing()

	var doc *enigma.Doc
	var err error

	doc, _ = global.GetActiveDoc(ctx)
	if doc != nil {
		// There is an already opened doc!
		log.Verboseln("App with name: " + appName + " and id: " + appID + "(reconnected)")
	} else {
		doc, err = global.OpenDoc(ctx, appID, "", "", "", noData)
		if doc != nil {
			if noData {
				log.Verboseln("Opened app with name: " + appName + " and id: " + appID + " without data")
			} else {
				log.Verboseln("Opened app with name: " + appName + " and id: " + appID)
			}
		} else if createAppIfMissing {
			var success bool
			success, appID, err = global.CreateApp(ctx, appName, "")
			if err != nil {
				log.Fatalf("could not create app with name '%s': %s\n", appName, err)
			}
			if !success {
				log.Fatalf("could not create app with name '%s'\n", appName)
			}
			// Write app id to config
			SetAppIDToKnownApps(engineHost, appName, appID, false)
			doc, err = global.OpenDoc(ctx, appID, "", "", "", noData)
			if err != nil {
				log.Fatalf("could not do open app with ID '%s': %s\n", appID, err)
			}
			if doc != nil {
				log.Verboseln("App with name: " + appName + " and id: " + appID + "(new)")
			}
		} else {
			log.Fatalf("could not open app with ID '%s': %s\n", appID, err)
		}
	}

	return doc
}

func logConnectError(err error, engine string) {
	msg := fmt.Sprintf("could not connect to engine on %s\nDetails: %s\n", engine, err)
	if strings.Contains(err.Error(), "401") {
		msg += fmt.Sprintln("This probably means that you have provided either incorrect or no authorization credentials.")
		msg += fmt.Sprintln("Check that the headers specified are correct.")
	} else if strings.Contains(err.Error(), "x509") {
		msg += fmt.Sprintln("This probably means that you have certificates that are not signed properly.")
		msg += fmt.Sprintln("If you have a self signed certificate then the flag '--insecure' might solve the problem.")
	} else {
		msg += fmt.Sprintln("This probably means that there is no engine running on the specified url.")
		msg += fmt.Sprintln("Check that the engine is up and that the url specified is correct.")
	}
	log.Fatalln(msg)
}

func GetGlobal(ctx context.Context, settings EngineWebSocketSettings) *enigma.Global {
	engineURL := settings.WebSocketEngineURL()
	fmt.Println("ENGINEURL.:..", engineURL)
	log.Verboseln("Engine: " + engineURL)

	headers := settings.Headers()
	if headers.Get("X-Qlik-Session") == "" {
		sessionID := getSessionID(settings.App(), settings.NoData(), settings.Ttl())
		headers.Set("X-Qlik-Session", sessionID)
	}
	log.Verboseln("SessionId " + headers.Get("X-Qlik-Session"))

	var dialer = enigma.Dialer{}
	dialer.TLSClientConfig = settings.TlsConfig()

	if log.Traffic {
		dialer.TrafficLogger = log.TrafficLogger{}
	}

	global, err := dialer.Dial(ctx, engineURL, settings.Headers())
	if err != nil {
		logConnectError(err, engineURL)
	}
	sessionMessages := global.SessionMessageChannel()
	err = waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		log.Fatalln("could not connect to engine: ", err)
	}
	go printSessionMessagesIfInVerboseMode(sessionMessages)
	return global
}

func waitForOnConnectedMessage(sessionMessages chan enigma.SessionMessage) error {
	for sessionEvent := range sessionMessages {
		log.Verboseln(sessionEvent.Topic + " " + string(sessionEvent.Content))
		if sessionEvent.Topic == "OnConnected" {
			var parsedEvent map[string]string
			err := json.Unmarshal(sessionEvent.Content, &parsedEvent)
			if err != nil {
				log.Fatalln("could not parse response from engine: ", err)
			}
			if parsedEvent["qSessionState"] == "SESSION_CREATED" || parsedEvent["qSessionState"] == "SESSION_ATTACHED" {
				return nil
			}
			return errors.New(parsedEvent["qSessionState"])
		}
	}
	return errors.New("session closed before reciving OnConnected message")
}

func printSessionMessagesIfInVerboseMode(sessionMessages chan enigma.SessionMessage) {
	for sessionEvent := range sessionMessages {
		log.Verboseln(sessionEvent.Topic + " " + string(sessionEvent.Content))
	}
}

func getSessionID(appID string, noData bool, ttl string) string {
	// If no-data or ttl flag is used the user should not get the default session id

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalln("unexpected error when retrieving current user: ", err)
	}
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatalln("unexpected error when retrieving hostname: ", err)
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte("corectl-" + currentUser.Username + "-" + hostName + "-" + appID + "-" + ttl + "-" + strconv.FormatBool(noData)))
	return sessionID
}
