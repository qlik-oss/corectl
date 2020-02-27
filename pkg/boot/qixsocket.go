package boot

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/user"
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
	WebSocketEngineURL() string
	AppIdMappingNamespace() string
}

func GetApp(ctx context.Context, global *enigma.Global, settings EngineWebSocketSettings, createAppIfMissing bool) *enigma.Doc {
	appName := settings.App()

	if appName == "" {
		log.Fatalln("no app specified")
	}

	appID := settings.AppId()
	noData := settings.NoData()

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
			SetAppIDToKnownApps(settings.AppIdMappingNamespace(), appName, appID, false)
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
