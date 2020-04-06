package boot

import (
	"context"
	"fmt"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/pkg/rest"
	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/cobra"
)

func NewCommunicator(ccmd *cobra.Command) *Communicator {
	commonSettings := Boot(ccmd)
	return &Communicator{CommonSettings: commonSettings}
}

// Boot reads all the configuration for a command and sets
// the log verbosity for the logger.
func Boot(ccmd *cobra.Command) *CommonSettings {
	cfg := dynconf.ReadSettings(ccmd)
	commonSettings := NewCommonSettings(cfg)
	log.Init(commonSettings.PrintMode())
	return commonSettings
}

type Communicator struct {
	*CommonSettings
}

func (c *Communicator) RestCaller() *rest.RestCaller {
	return &rest.RestCaller{RestCallerSettings: c.CommonSettings}
}

func (c *Communicator) OpenAppSocket(createAppIfMissing bool) (context.Context, *enigma.Global, *enigma.Doc, *CommonSettings) {
	if c.IsSenseForKubernetes() {
		c.ExplicitlyTranslatedAppId = c.RestCaller().TranslateAppNameToId(c.App())
	}
	global := GetGlobal(context.Background(), c.CommonSettings)
	app := GetApp(context.Background(), global, c.CommonSettings, createAppIfMissing)
	return context.Background(), global, app, c.CommonSettings
}
func (c *Communicator) OpenGlobalSocket() (context.Context, *enigma.Global, *CommonSettings) {
	if c.IsSenseForKubernetes() {
		log.Fatalln("Not implemented for Sense yet")
	}
	global := GetGlobal(context.Background(), c.CommonSettings)
	return context.Background(), global, c.CommonSettings
}

//AppExists returns wether or not an app exists along with any eventual error from engine.
func (c *Communicator) AppExists() (bool, error) {
	if c.IsSenseForKubernetes() {
		log.Fatalln("Not implemented for Sense yet")
	}
	global := GetGlobal(context.Background(), c.CommonSettings)
	_, err := global.GetAppEntry(context.Background(), c.CommonSettings.AppId())
	if err != nil {
		return false, fmt.Errorf("could not find any app by ID '%s': %s", c.CommonSettings.AppId(), err)
	}
	return true, nil
}

//DeleteApp removes the specified app from the engine.
func (c *Communicator) DeleteApp(appName string) {
	if c.IsSenseForKubernetes() {
		log.Fatalln("Not implemented for Sense yet")
	}
	c.CommonSettings.OverrideSetting("app", appName)
	global := GetGlobal(context.Background(), c.CommonSettings)
	appID := c.CommonSettings.AppId()
	succ, err := global.DeleteApp(context.Background(), appID)
	if err != nil {
		log.Fatalf("could not delete app with name '%s' and ID '%s': %s\n", appName, appID, err)
	} else if !succ {
		log.Fatalf("could not delete app with name '%s' and ID '%s'\n", appName, appID)
	}
	SetAppIDToKnownApps(c.CommonSettings.AppIdMappingNamespace(), appName, appID, true)
}
