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
		Logger.Info("Could not connect to the default engine on http://localhost:9076")
		Logger.Info("Specify where the engine is running using the --engine parameter or in your config file.")
		Logger.Fatal("Error details: ", err)
	} else {
		Logger.Info("Could not connect to engine on " + engine + ".")
		Logger.Info("Please check the --engine parameter or your config file.")
		Logger.Fatal("Error details: ", err)
	}
}

func connectToEngine(ctx context.Context, engine string, appID string, ttl string, headers http.Header) *enigma.Global {
	engineURL := buildWebSocketURL(engine, ttl)
	Logger.Debug("Engine: " + engineURL)

	if headers.Get("X-Qlik-Session") == "" {
		sessionID := getSessionID(appID)
		headers.Set("X-Qlik-Session", sessionID)
	}
	Logger.Debug("SessionId: " + headers.Get("X-Qlik-Session"))

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
		Logger.Fatal(err)
	} else if !succ {
		Logger.Fatal("Failed to delete app " + appID)
	}
}

// PrepareEngineState makes sure that the app idenfied by the supplied parameters is created or opened or reconnected to
// depending on the state. The TTL feature is used to keep the app session loaded to improve performance.
func PrepareEngineState(ctx context.Context, engine string, appID string, ttl string, headers http.Header, createAppIfMissing bool) *State {
	Logger.Debug("---------- Connecting to app ----------")
	global := connectToEngine(ctx, engine, appID, ttl, headers)
	if appID == "" {
		Logger.Info("No app specified, using session app.")
	}
	sessionMessages := global.SessionMessageChannel()
	err := waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		Logger.Fatal("Failed to connect to engine with error message: ", err)
	}
	go printSessionMessagesIfInVerboseMode(sessionMessages)
	doc, err := global.GetActiveDoc(ctx)
	if doc != nil {
		// There is an already opened doc!
		if appID != "" {
			Logger.Debug("App: " + appID + "(reconnected)")
		} else {
			Logger.Debug("Session app (reconnected)")
		}
	} else {
		if appID == "" {
			doc, err = global.CreateSessionApp(ctx)
			if doc != nil {
				Logger.Debug("Session app (new)")
			} else {
				Logger.Fatal(err)
			}
		} else {
			doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
			if doc != nil {
				Logger.Debug("App:  " + appID + "(opened)")
			} else if createAppIfMissing {
				_, _, err = global.CreateApp(ctx, appID, "")
				if err != nil {
					Logger.Fatal(err)
				}
				doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
				if err != nil {
					Logger.Fatal(err)
				}
				if doc != nil {
					Logger.Debug("Document: " + appID + "(new)")
				}
			} else {
				Logger.Fatal(err)
			}
		}
	}

	metaURL := buildMetadataURL(engine, appID)
	Logger.Debug("Meta: " + metaURL)

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
		Logger.Debug(sessionEvent.Topic + " " + string(sessionEvent.Content))
		if sessionEvent.Topic == "OnConnected" {
			var parsedEvent map[string]string
			err := json.Unmarshal(sessionEvent.Content, &parsedEvent)
			if err != nil {
				Logger.Fatal(err)
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
		Logger.Debug(sessionEvent.Topic + " " + string(sessionEvent.Content))
	}
}

// PrepareEngineStateWithoutApp creates a connection to the engine with no dependency to any app.
func PrepareEngineStateWithoutApp(ctx context.Context, engine string, ttl string, headers http.Header) *State {
	Logger.Debug("---------- Connecting to engine ----------")

	engineURL := buildWebSocketURL(engine, ttl)

	Logger.Debug("Engine: " + engineURL)

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
		Logger.Fatal("Failed to connect to engine with error message: ", err)
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
		Logger.Fatal(err)
	}
	hostName, err := os.Hostname()
	if err != nil {
		Logger.Fatal(err)
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte("Corectl-" + currentUser.Username + "-" + hostName + "-" + appID))
	return sessionID
}
