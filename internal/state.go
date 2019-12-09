package internal

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
	"github.com/spf13/viper"
)

// State contains all needed info about the current app including a go context to use when communicating with the engine.
type State struct {
	Doc     *enigma.Doc
	Ctx     context.Context
	Global  *enigma.Global
	AppName string
	AppID   string
	Verbose bool
}

func logConnectError(err error, engine string) {
	msg := fmt.Sprintf("could not connect to engine on %s\nDetails: %s\n", engine, err)
	if strings.Contains(err.Error(), "401") {
		msg += fmt.Sprintln("This probably means that you have provided either incorrect or no authorization credentials.")
		msg += fmt.Sprintln("Check that the headers specified are correct.")
	} else {
		msg += fmt.Sprintln("This probably means that there is no engine running on the specified url.")
		msg += fmt.Sprintln("Check that the engine is up and that the url specified is correct.")
	}
	log.Fatalln(msg)
}

func connectToEngine(ctx context.Context, appName, ttl string, headers http.Header, tlsClientConfig *tls.Config) *enigma.Global {
	engineURL := buildWebSocketURL(ttl)
	log.Verboseln("Engine: " + engineURL)

	if headers.Get("X-Qlik-Session") == "" {
		sessionID := getSessionID(appName)
		headers.Set("X-Qlik-Session", sessionID)
	}
	log.Verboseln("SessionId " + headers.Get("X-Qlik-Session"))

	var dialer = enigma.Dialer{}
	dialer.TLSClientConfig = tlsClientConfig

	if log.Traffic {
		dialer.TrafficLogger = log.TrafficLogger{}
	}

	global, err := dialer.Dial(ctx, engineURL, headers)
	if err != nil {
		logConnectError(err, engineURL)
	}
	return global
}

//AppExists returns wether or not an app exists along with any eventual error from engine.
func AppExists(ctx context.Context, engine string, appName string, headers http.Header, tlsClientConfig *tls.Config) (bool, error) {
	global := PrepareEngineState(ctx, headers, tlsClientConfig, false, true).Global
	appID, _ := applyNameToIDTransformation(appName)
	_, err := global.GetAppEntry(ctx, appID)
	if err != nil {
		return false, fmt.Errorf("could not find any app by ID '%s': %s", appID, err)
	}
	return true, nil
}

//DeleteApp removes the specified app from the engine.
func DeleteApp(ctx context.Context, engine string, appName string, headers http.Header, tlsClientConfig *tls.Config) {
	global := PrepareEngineState(ctx, headers, tlsClientConfig, false, true).Global
	appID, _ := applyNameToIDTransformation(appName)
	succ, err := global.DeleteApp(ctx, appID)
	if err != nil {
		log.Fatalf("could not delete app with name '%s' and ID '%s': %s\n", appName, appID, err)
	} else if !succ {
		log.Fatalf("could not delete app with name '%s' and ID '%s'\n", appName, appID)
	}
	SetAppIDToKnownApps(appName, appID, true)
}

// PrepareEngineState connects to engine with or without an app. It returns a *State
// which contains at least *enigma.Global (connection to engine used to open apps) and a context.Context.
// If called with app (withoutApp = false) it will also return *enigma.Doc (the representation of a qlik-app)
// as well as app name and ID.
//
// Any ttl supplied (through viper) specifies how long the engine should keep the session alive which affects
// performance. (It is cheaper to reattach to a pre-existing session, performance-wise.)
func PrepareEngineState(ctx context.Context, headers http.Header, tlsClientConfig *tls.Config, createAppIfMissing, withoutApp bool) *State {
	engine := viper.GetString("engine")
	appName := viper.GetString("app")
	ttl := viper.GetString("ttl")
	noData := viper.GetBool("no-data")

	var doc *enigma.Doc
	var appID string

	// If no app was supplied but is needed, check the url.
	if appName == "" && !withoutApp {
		appName = TryParseAppFromURL(engine)

		if appName == "" {
			log.Fatalln("no app specified")
		}
	}

	log.Verboseln("---------- Connecting to engine ----------")
	global := connectToEngine(ctx, appName, ttl, headers, tlsClientConfig)
	sessionMessages := global.SessionMessageChannel()
	err := waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		log.Fatalln("could not connect to engine: ", err)
	}
	go printSessionMessagesIfInVerboseMode(sessionMessages)

	if !withoutApp {
		appID, _ = applyNameToIDTransformation(appName)
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
				SetAppIDToKnownApps(appName, appID, false)
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
	}

	return &State{
		Doc:     doc,
		Global:  global,
		AppName: appName,
		AppID:   appID,
		Ctx:     ctx,
	}
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

func getSessionID(appID string) string {
	// If no-data or ttl flag is used the user should not get the default session id
	noData := viper.GetBool("no-data")
	ttl := viper.GetString("ttl")

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
