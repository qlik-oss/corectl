package internal

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/qlik-oss/corectl/internal/log"
	"gopkg.in/yaml.v2"
)

var knownAppsFilePath = path.Join(userHomeDir(), ".corectl", "knownApps.yml")

// Fetch a matching app id from known apps for a specified app name
// If not found return the appName and found bool set to false
func applyNameToIDTransformation(appName string) (appID string, found bool) {
	apps := getKnownApps()

	if apps == nil {
		log.Verboseln("knownApps yaml file not found")
		return appName, false
	}

	engineURL := GetEngineURL()
	host := engineURL.Host

	if id, exists := apps[host][appName]; exists {
		log.Verboseln("Found id: " + id + " for app with name: " + appName + " @" + host)
		return id, true
	}

	log.Verboseln("No known id for app with name: " + appName)
	return appName, false
}

// Get map of known apps
func getKnownApps() map[string]map[string]string {
	var knownApps = map[string]map[string]string{}
	yamlFile, err := ioutil.ReadFile(knownAppsFilePath)
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(yamlFile, &knownApps)
	if err != nil {
		log.Fatalf("could not parse content of knownApps yaml '%s': %s\n", yamlFile, err)
	}

	return knownApps
}

// SetAppIDToKnownApps adds an app or removes an app from known apps
func SetAppIDToKnownApps(appName string, appID string, remove bool) {

	createKnownAppsFileIfNotExist()
	apps := getKnownApps()

	engineURL := GetEngineURL()
	host := engineURL.Host

	// Either remove or add an entry
	if remove {
		if _, exists := apps[host][appName]; exists {
			delete(apps[host], appName)
			log.Verboseln("Removed app with name: " + appName + " and id: " + appID + " @" + host + " from known apps")
		}
	} else {
		if apps[host] == nil {
			apps[host] = map[string]string{}
		}
		apps[host][appName] = appID
		log.Verboseln("Added app with name: " + appName + " and id: " + appID + " @" + host + " to known apps")
	}

	// Write to knownApps.yml
	out, _ := yaml.Marshal(apps)

	if err := ioutil.WriteFile(knownAppsFilePath, out, 0644); err != nil {
		log.Fatalf("could not write to '%s': %s\n", knownAppsFilePath, err)
	}
}

// Create a knownApps.yml if one does not exist
func createKnownAppsFileIfNotExist() {
	if _, err := os.Stat(knownAppsFilePath); os.IsNotExist(err) {

		// Create .corectl folder in home directory
		corectlDir := path.Join(userHomeDir(), ".corectl")
		if _, err := os.Stat(corectlDir); os.IsNotExist(err) {
			err = os.Mkdir(corectlDir, os.ModePerm)
			if err != nil {
				log.Fatalln("could not create .corectl folder in home directory: ", err)
			}
		}

		// Create knownApps.yml in .corectl folder
		_, err := os.Create(knownAppsFilePath)
		if err != nil {
			log.Fatalf("could not create %s: %s\n", knownAppsFilePath, err)
		}

		log.Verboseln("Created ~/.corectl/knownApps.yml for storage of app ids")
	}
}

// Get the user home directory dependent on OS
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
