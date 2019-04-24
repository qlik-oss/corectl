package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"os"
	"os/user"
	"strconv"
	"strings"

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
	MetaURL string
	Verbose bool
}

func logConnectError(err error, engine string) {

	if engine == "" {
		fmt.Println("Could not connect to the default engine on http://localhost:9076")
		fmt.Println("Specify where the engine is running using the --engine parameter or in your config file.")
		fmt.Println("Error details: ", err)
	} else {
		fmt.Println("Could not connect to engine on " + engine + ".")
		fmt.Println("Please check the --engine parameter or your config file.")
		fmt.Println("Error details: ", err)
	}
	os.Exit(1)
}

func connectToEngine(ctx context.Context, engine string, appName string, headers http.Header) *enigma.Global {
	engineURL := BuildWebSocketURL(engine, false)
	LogVerbose("Engine: " + engineURL)

	if headers.Get("X-Qlik-Session") == "" {
		sessionID := getSessionID(appName)
		LogVerbose("SessionId: " + sessionID)
		headers.Set("X-Qlik-Session", sessionID)
	}
	LogVerbose("SessionId " + headers.Get("X-Qlik-Session"))

	var dialer enigma.Dialer

	if LogTraffic {
		dialer = enigma.Dialer{TrafficLogger: TrafficLogger{}}
	} else {
		dialer = enigma.Dialer{}
	}

	global, err := dialer.Dial(ctx, engineURL, headers)
	if err != nil {
		logConnectError(err, engine)
	}
	return global
}

//AppExists returns wether or not an app exists
func AppExists(ctx context.Context, engine string, appName string, headers http.Header) bool {
	global := connectToEngine(ctx, engine, appName, headers)
	appID, _ := applyNameToIDTransformation(engine, appName)
	_, err := global.GetAppEntry(ctx, appID)
	return err == nil
}

//DeleteApp removes the specified app from the engine.
func DeleteApp(ctx context.Context, engine string, appName string, headers http.Header) {
	global := connectToEngine(ctx, engine, appName, headers)
	appID, _ := applyNameToIDTransformation(engine, appName)
	succ, err := global.DeleteApp(ctx, appID)
	if err != nil {
		FatalError(err)
	} else if !succ {
		FatalError("Failed to delete app with name: " + appName + " and ID: " + appID)
	}
	setAppIDToKnownApps(engine, appName, appID, true)
}

// PrepareEngineState makes sure that the app idenfied by the supplied parameters is created or opened or reconnected to
// depending on the state. The TTL feature is used to keep the app session loaded to improve performance.
func PrepareEngineState(ctx context.Context, headers http.Header, createAppIfMissing bool) *State {
	engine := viper.GetString("engine")
	appName := viper.GetString("app")
	noData := viper.GetBool("no-data")
	bashMode := viper.GetBool("bash")
	var appID string

	LogVerbose("---------- Connecting to app ----------")
	global := connectToEngine(ctx, engine, appName, headers)
	if appName == "" && !bashMode {
		fmt.Println("No app specified, using session app.")
	}
	sessionMessages := global.SessionMessageChannel()
	err := waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		FatalError("Failed to connect to engine with error message: ", err)
	}
	go printSessionMessagesIfInVerboseMode(sessionMessages)
	doc, err := global.GetActiveDoc(ctx)
	if doc != nil {
		// There is an already opened doc!
		if appName != "" {
			appID, _ = applyNameToIDTransformation(engine, appName)
			LogVerbose("App with name: " + appName + " and id: " + appID + "(reconnected)")
		} else {
			LogVerbose("Session app (reconnected)")
		}
	} else {
		if appName == "" {
			doc, err = global.CreateSessionApp(ctx)
			if doc != nil {
				LogVerbose("Session app (new)")
			} else {
				FatalError(err)
			}
		} else {
			appID, _ = applyNameToIDTransformation(engine, appName)
			doc, err = global.OpenDoc(ctx, appID, "", "", "", noData)
			if doc != nil {
				if noData {
					LogVerbose("Opened app with name: " + appName + " and id: " + appID + " without data")
				} else {
					LogVerbose("Opened app with name: " + appName + " and id: " + appID)
				}
			} else if createAppIfMissing {
				success, appID, err := global.CreateApp(ctx, appName, "")
				if err != nil {
					FatalError(err)
				}
				if !success {
					FatalError("Failed to create app with name: " + appName)
				}
				// Write app id to config
				setAppIDToKnownApps(engine, appName, appID, false)
				doc, err = global.OpenDoc(ctx, appID, "", "", "", noData)
				if err != nil {
					FatalError(err)
				}
				if doc != nil {
					LogVerbose("App with name: " + appName + " and id: " + appID + "(new)")
				}
			} else {
				FatalError(err)
			}
		}
	}

	metaURL := buildMetadataURL(engine, appID)
	LogVerbose("Meta: " + metaURL)

	return &State{
		Doc:     doc,
		Global:  global,
		AppName: appName,
		AppID:   appID,
		Ctx:     ctx,
		MetaURL: metaURL,
	}
}

func waitForOnConnectedMessage(sessionMessages chan enigma.SessionMessage) error {
	for sessionEvent := range sessionMessages {
		LogVerbose(sessionEvent.Topic + " " + string(sessionEvent.Content))
		if sessionEvent.Topic == "OnConnected" {
			var parsedEvent map[string]string
			err := json.Unmarshal(sessionEvent.Content, &parsedEvent)
			if err != nil {
				FatalError(err)
			}
			if parsedEvent["qSessionState"] == "SESSION_CREATED" || parsedEvent["qSessionState"] == "SESSION_ATTACHED" {
				return nil
			}
			return errors.New(parsedEvent["qSessionState"])
		}
	}
	return errors.New("Session closed before reciving OnConnected message.")
}

func printSessionMessagesIfInVerboseMode(sessionMessages chan enigma.SessionMessage) {
	for sessionEvent := range sessionMessages {
		LogVerbose(sessionEvent.Topic + " " + string(sessionEvent.Content))
	}
}

// PrepareEngineStateWithoutApp creates a connection to the engine with no dependency to any app.
func PrepareEngineStateWithoutApp(ctx context.Context, headers http.Header) *State {
	engine := viper.GetString("engine")

	LogVerbose("---------- Connecting to engine ----------")

	engineURL := BuildWebSocketURL(engine, false)

	LogVerbose("Engine: " + engineURL)

	var dialer enigma.Dialer

	if LogTraffic {
		dialer = enigma.Dialer{TrafficLogger: TrafficLogger{}}
	} else {
		dialer = enigma.Dialer{}
	}

	global, err := dialer.Dial(ctx, engineURL, headers)

	if err != nil {
		logConnectError(err, engine)
	}
	sessionMessages := global.SessionMessageChannel()
	err = waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		FatalError("Failed to connect to engine with error message: ", err)
	}
	go printSessionMessagesIfInVerboseMode(sessionMessages)

	return &State{
		Doc:     nil,
		Global:  global,
		AppName: "",
		AppID:   "",
		Ctx:     ctx,
		MetaURL: "",
	}
}

// BuildWebSocketURL modifies the websocket url if incomplete
func BuildWebSocketURL(engine string, baseURLOnly bool) string {
	if strings.HasPrefix(engine, "wss://") || strings.HasPrefix(engine, "ws://") {
		return engine
	} else if baseURLOnly {
		return "ws://" + engine
	} else {
		ttl := viper.GetString("ttl")
		return "ws://" + engine + "/app/ttl/" + ttl
	}
}

func buildMetadataURL(engine string, appID string) string {
	if appID == "" {
		return ""
	}
	engine = BuildWebSocketURL(engine, true)
	engine = strings.Replace(engine, "wss://", "https://", -1)
	engine = strings.Replace(engine, "ws://", "http://", -1)
	url := fmt.Sprintf("%s/v1/apps/%s/data/metadata", engine, neturl.QueryEscape(appID))
	return url
}

func getSessionID(appID string) string {
	// If no-data or ttl flag is used the user should not get the default session id
	noData := viper.GetBool("no-data")
	ttl := viper.GetString("ttl")

	currentUser, err := user.Current()
	if err != nil {
		FatalError(err)
	}
	hostName, err := os.Hostname()
	if err != nil {
		FatalError(err)
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte("corectl-" + currentUser.Username + "-" + hostName + "-" + appID + "-" + ttl + "-" + strconv.FormatBool(noData)))
	return sessionID
}
