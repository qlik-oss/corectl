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
	"strings"

	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

// State contains all needed info about the current app including a go context to use when communicating with the engine.
type State struct {
	Doc     *enigma.Doc
	Ctx     context.Context
	Global  *enigma.Global
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

func connectToEngine(ctx context.Context, engine string, appID string, ttl string, headers http.Header) *enigma.Global {
	engineURL := buildWebSocketURL(engine, ttl)
	LogVerbose("Engine: " + engineURL)

	if headers.Get("X-Qlik-Session") == "" {
		sessionID := getSessionID(appID)
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

//DeleteApp removes the specified app from the engine.
func DeleteApp(ctx context.Context, engine string, appID string, ttl string, headers http.Header) {
	global := connectToEngine(ctx, engine, appID, ttl, headers)
	succ, err := global.DeleteApp(ctx, appID)
	if err != nil {
		FatalError(err)
	} else if !succ {
		FatalError("Failed to delete app " + appID)
	}
}

// PrepareEngineState makes sure that the app idenfied by the supplied parameters is created or opened or reconnected to
// depending on the state. The TTL feature is used to keep the app session loaded to improve performance.
func PrepareEngineState(ctx context.Context, headers http.Header, createAppIfMissing bool) *State {
	engine := viper.GetString("engine")
	appID := viper.GetString("app")
	ttl := viper.GetString("ttl")
	noData := viper.GetBool("no-data")

	LogVerbose("---------- Connecting to app ----------")
	global := connectToEngine(ctx, engine, appID, ttl, headers)
	if appID == "" {
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
		if appID != "" {
			LogVerbose("App: " + appID + "(reconnected)")
		} else {
			LogVerbose("Session app (reconnected)")
		}
	} else {
		if appID == "" {
			doc, err = global.CreateSessionApp(ctx)
			if doc != nil {
				LogVerbose("Session app (new)")
			} else {
				FatalError(err)
			}
		} else {
			doc, err = global.OpenDoc(ctx, appID, "", "", "", noData)
			if doc != nil {
				if noData {
					LogVerbose("Opened app with id " + appID + " without data")
				} else {
					LogVerbose("Opened app with id " + appID)
				}
			} else if createAppIfMissing {
				_, _, err = global.CreateApp(ctx, appID, "")
				if err != nil {
					FatalError(err)
				}
				doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
				if err != nil {
					FatalError(err)
				}
				if doc != nil {
					LogVerbose("Document: " + appID + "(new)")
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
	return errors.New("Session closed before reciving OnConnected message")
}

func printSessionMessagesIfInVerboseMode(sessionMessages chan enigma.SessionMessage) {
	for sessionEvent := range sessionMessages {
		LogVerbose(sessionEvent.Topic + " " + string(sessionEvent.Content))
	}
}

// PrepareEngineStateWithoutApp creates a connection to the engine with no dependency to any app.
func PrepareEngineStateWithoutApp(ctx context.Context, headers http.Header) *State {
	engine := viper.GetString("engine")
	ttl := viper.GetString("ttl")

	LogVerbose("---------- Connecting to engine ----------")

	engineURL := buildWebSocketURL(engine, ttl)

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
		AppID:   "",
		Ctx:     ctx,
		MetaURL: "",
	}
}

//TidyUpEngineURL tidies up an engine url fragment and returns a complete url.
func TidyUpEngineURL(engine string) string {
	if engine == "" {
		engine = "localhost:9076"
	}
	var url string
	if strings.HasPrefix(engine, "wss://") {
		url = engine
	} else if strings.HasPrefix(engine, "ws://") {
		url = engine
	} else {
		url = "ws://" + engine
	}
	if len(strings.Split(url, ":")) == 2 {
		url += ":9076"
	}
	return url
}

func buildWebSocketURL(engine string, ttl string) string {
	engine = TidyUpEngineURL(engine)
	return engine + "/app/engineData/ttl/" + ttl
}

func buildMetadataURL(engine string, appID string) string {
	if appID == "" {
		return ""
	}
	engine = TidyUpEngineURL(engine)
	engine = strings.Replace(engine, "wss://", "https://", -1)
	engine = strings.Replace(engine, "ws://", "http://", -1)
	url := fmt.Sprintf("%s/v1/apps/%s/data/metadata", engine, neturl.QueryEscape(appID))
	return url
}

func getSessionID(appID string) string {
	currentUser, err := user.Current()
	if err != nil {
		FatalError(err)
	}
	hostName, err := os.Hostname()
	if err != nil {
		FatalError(err)
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte("Corectl-" + currentUser.Username + "-" + hostName + "-" + appID))
	return sessionID
}

// FatalError prints the supplied message and exists the process with code 1
func FatalError(fatalMessage ...interface{}) {
	fmt.Println(fatalMessage...)
	os.Exit(1)
}
