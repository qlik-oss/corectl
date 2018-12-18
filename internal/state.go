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

	fmt.Println("before connected!")
	global, err := enigma.Dialer{}.Dial(ctx, engineURL, headers)
	if err != nil {
		logConnectError(err, engine)
	}
	fmt.Println("I connected!")
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
func PrepareEngineState(ctx context.Context, engine string, appID string, ttl string, headers http.Header, createAppIfMissing bool) *State {
	LogVerbose("---------- Connecting to app ----------")
	global := connectToEngine(ctx, engine, appID, ttl, headers)
	if appID == "" {
		fmt.Println("No app specified, using session app.")
	}
	go func() {
		for x := range global.SessionMessageChannel() {
			if x.Topic != "OnConnected" {
				fmt.Println(x.Topic, string(x.Content))
			}
		}
	}()
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
			doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
			if doc != nil {
				LogVerbose("App:  " + appID + "(opened)")
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

//Låsa tills onConnected?
//Allmänt printa på very verbose?
func printMessage(global *enigma.Global) {
	sessionMessages := global.SessionMessageChannel()
	for sessionEvent := range sessionMessages {
		fmt.Println("TOPIC: ", sessionEvent.Topic)
		//fmt.Println("Content: ", sessionEvent.Content)
		var parsed map[string]string
		err := json.Unmarshal(sessionEvent.Content, &parsed)
		fmt.Println("???")
		if err != nil {
			fmt.Println("Error parsing ", err)
			// for key := range parsed.key()

		} else {
			fmt.Println("dags att visa upp vad jag har: ")
			for k, v := range parsed {
				fmt.Printf("key[%s] value[%s]\n", k, v)
			}
		}
	}
}

func waitForOnConnectedMessage(sessionMessages chan enigma.SessionMessage) error {
	for sessionEvent := range sessionMessages {
		if sessionEvent.Topic == "OnConnected" {
			var parsed map[string]string
			err := json.Unmarshal(sessionEvent.Content, &parsed)
			if err != nil {
				FatalError(err)
			}
			if parsed["qSessionState"] == "SESSION_CREATED" || parsed["qSessionState"] == "SESSION_ATTACHED" {
				return nil
			}
			return errors.New(parsed["qSessionState"])

			// } else if parsed["qSessionState"] == "SESSION_ERROR_NO_LICENSE" {
			// 	FatalError(errors.New("Failed to connect to engine because of missing license"))
			// } else {
			// 	for k, v := range parsed {
			// 		fmt.Printf("key[%s] value[%s]\n", k, v)
			// 	}
			// }
		}
	}
	return errors.New("This should never happen")
}

// PrepareEngineStateWithoutApp creates a connection to the engine with no dependency to any app.
func PrepareEngineStateWithoutApp(ctx context.Context, engine string, ttl string, headers http.Header) *State {
	LogVerbose("---------- Connecting to engine ----------")

	engineURL := buildWebSocketURL(engine, ttl)

	LogVerbose("Engine: " + engineURL)
	fmt.Println("Before connection")
	global, err := enigma.Dialer{}.Dial(ctx, engineURL, headers)

	if err != nil {
		logConnectError(err, engine)
	}
	sessionMessages := global.SessionMessageChannel()
	err = waitForOnConnectedMessage(sessionMessages)
	if err != nil {
		FatalError("Failed to connect to engine with error message: ", err)
	}
	// for sessionEvent := range sessionMessages {
	// 	// var parsed map[string]string
	// 	// err := json.Unmarshal(sessionEvent.Content, &parsed)
	// 	// if err != nil {
	// 	// 	FatalError(err)
	// 	// }

	// 	//bryt ut till funktion som returnerar eror eller nil (Ha sen en annan funktion som körs kontiunerligt)
	// 	if sessionEvent.Topic == "OnConnected" {
	// 		var parsed map[string]string
	// 		err := json.Unmarshal(sessionEvent.Content, &parsed)
	// 		if err != nil {
	// 			FatalError(err)
	// 		}
	// 		if parsed["qSessionState"] == "SESSION_CREATED" || parsed["qSessionState"] == "SESSION_ATTACHED" {
	// 			break
	// 		} else if parsed["qSessionState"] == "SESSION_ERROR_NO_LICENSE" {
	// 			FatalError(errors.New("Failed to connect to engine because of missing license"))
	// 		} else {
	// 			for k, v := range parsed {
	// 				fmt.Printf("key[%s] value[%s]\n", k, v)
	// 			}
	// 		}
	// 	}
	// }

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
