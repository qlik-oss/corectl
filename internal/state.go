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
	"regexp"
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

func connectToEngine(ctx context.Context, engine string, appName string, ttl string, headers http.Header) *enigma.Global {
	engineURL := buildWebSocketURL(engine, ttl)
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
	global := PrepareEngineStateWithoutApp(ctx, headers).Global
	appID, _ := applyNameToIDTransformation(engine, appName)
	_, err := global.GetAppEntry(ctx, appID)
	return err == nil
}

//DeleteApp removes the specified app from the engine.
func DeleteApp(ctx context.Context, engine string, appName string, headers http.Header) {
	global := PrepareEngineStateWithoutApp(ctx, headers).Global
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
	ttl := viper.GetString("ttl")
	noData := viper.GetBool("no-data")

	if appName == "" {
		// No app name provided, lets check if one exists in the url
		appName = tryParseAppFromURL(engine)
		if appName == "" {
			FatalError("Error: No app specified")
		}
	}

	appID, _ := applyNameToIDTransformation(engine, appName)

	LogVerbose("---------- Connecting to app ----------")
	global := connectToEngine(ctx, engine, appName, ttl, headers)
	sessionMessages := global.SessionMessageChannel()
	err := waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		FatalError("Failed to connect to engine with error message: ", err)
	}
	go printSessionMessagesIfInVerboseMode(sessionMessages)
	doc, err := global.GetActiveDoc(ctx)
	if doc != nil {
		// There is an already opened doc!
		LogVerbose("App with name: " + appName + " and id: " + appID + "(reconnected)")
	} else {
		doc, err = global.OpenDoc(ctx, appID, "", "", "", noData)
		if doc != nil {
			if noData {
				LogVerbose("Opened app with name: " + appName + " and id: " + appID + " without data")
			} else {
				LogVerbose("Opened app with name: " + appName + " and id: " + appID)
			}
		} else if createAppIfMissing {
			var success bool
			success, appID, err = global.CreateApp(ctx, appName, "")
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
	return errors.New("session closed before reciving OnConnected message")
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
		AppName: "",
		AppID:   "",
		Ctx:     ctx,
		MetaURL: "",
	}
}

//TidyUpEngineURL tidies up an engine url fragment and returns a complete url.
func TidyUpEngineURL(engine string) string {
	if strings.HasPrefix(engine, "wss://") || strings.HasPrefix(engine, "ws://") {
		return engine
	}
	return "ws://" + engine
}

func buildWebSocketURL(engine string, ttl string) string {
	engine = TidyUpEngineURL(engine)
	u, _ := neturl.Parse(engine)
	// Only modify the url if it does not contain a path or ends with a "/"
	if u.Path == "" && engine[len(engine)-1:] != "/" {
		engine = engine + "/app/engineData/ttl/" + ttl
	}
	return engine
}

func buildMetadataURL(engine string, appID string) string {
	engine = TidyUpEngineURL(engine)
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

// Function for parsing an url for an app identifier
func tryParseAppFromURL(engineURL string) string {
	// Find any string in the path succeeding "/app/", and excluding anything after "/"
	re, _ := regexp.Compile("/app/([^/]+)")
	values := re.FindStringSubmatch(engineURL)
	if len(values) > 0 {
		appName := values[1]
		LogVerbose("Found app in engine url: " + appName)
		return appName
	}
	return ""
}
