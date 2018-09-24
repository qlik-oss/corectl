package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
	"os"
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

// PrepareEngineState makes sure that the app idenfied by the supplied parameters is created or opened or reconnected to
// depending on the state. The TTL feature is used to keep the app session loaded to improve performance.
func PrepareEngineState(ctx context.Context, engine string, sessionID string, appID string, ttl string) *State {
	LogVerbose("---------- Connecting to app ----------")

	engineURL := buildWebSocketURL(engine, ttl)
	var doc *enigma.Doc

	LogVerbose("Engine: " + engineURL)

	LogVerbose("SessionId: " + sessionID)
	headers := make(http.Header, 1)
	headers.Set("X-Qlik-Session", sessionID)
	global, err := enigma.Dialer{}.Dial(ctx, engineURL, headers)
	if err != nil {
		fmt.Println("Could not connect to engine:"+engine, err)
		os.Exit(1)
	}

	go func() {
		for x := range global.SessionMessageChannel() {
			if x.Topic != "OnConnected" {
				fmt.Println(x.Topic, string(x.Content))
			}
		}
	}()
	doc, err = global.GetActiveDoc(ctx)
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
				log.Fatalln(err)
			}
		} else {
			doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
			if doc != nil {
				LogVerbose("App:  " + appID + "(opened)")
			} else {
				_, _, err = global.CreateApp(ctx, appID, "")
				if err != nil {
					log.Fatalln(err)
				}
				doc, err = global.OpenDoc(ctx, appID, "", "", "", false)
				if err != nil {
					log.Fatalln(err)
				}
				if doc != nil {
					LogVerbose("Document: " + appID + "(new)")
				}
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

// PrepareEngineStateWithoutApp creates a connection to the engine with no dependency to any app.
func PrepareEngineStateWithoutApp(ctx context.Context, engine string, ttl string) *State {
	LogVerbose("---------- Connecting to engine ----------")

	engineURL := buildWebSocketURL(engine, ttl)

	LogVerbose("Engine: " + engineURL)

	global, err := enigma.Dialer{}.Dial(ctx, engineURL, nil)
	if err != nil {
		fmt.Println("Could not connect to engine:"+engine, err)
		os.Exit(1)
	}

	return &State{
		Doc:     nil,
		Global:  global,
		AppID:   "",
		Ctx:     ctx,
		MetaURL: "",
	}
}

func tidyUpEngineURL(engine string) string {
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
	engine = tidyUpEngineURL(engine)
	return engine + "/app/engineData/ttl/" + ttl
}

func buildMetadataURL(engine string, appID string) string {
	if appID == "" {
		return ""
	}
	engine = tidyUpEngineURL(engine)
	engine = strings.Replace(engine, "wss://", "https://", -1)
	engine = strings.Replace(engine, "ws://", "http://", -1)
	url := fmt.Sprintf("%s/v1/apps/%s/data/metadata", engine, neturl.QueryEscape(appID))
	return url
}
