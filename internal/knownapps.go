package internal

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"gopkg.in/yaml.v2"
)

var knownAppsFilePath = path.Join(userHomeDir(), ".corectl", "knownApps.yml")

// Fetch a matching app id from known apps for a specified app name
// If not found return the appName and found bool set to false
func applyNameToIDTransformation(appName string) (appID string, found bool) {
	apps := getKnownApps()

	if apps == nil {
		LogVerbose("knownApps yaml file not found")
		return appName, false
	}

	engineURL := GetEngineURL()
	host := engineURL.Host

	if id, exists := apps[host][appName]; exists {
		LogVerbose("Found id: " + id + " for app with name: " + appName + " @" + host)
		return id, true
	}

	LogVerbose("No known id for app with name: " + appName)
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
		FatalErrorf("could not parse content of knownApps yaml '%s': %s", yamlFile, err)
	}

	return knownApps
}

// Add an app or remove an app from known apps
func SetAppIDToKnownApps(appName string, appID string, remove bool) {

	createKnownAppsFileIfNotExist()
	apps := getKnownApps()

	engineURL := GetEngineURL()
	host := engineURL.Host

	// Either remove or add an entry
	if remove {
		if _, exists := apps[host][appName]; exists {
			delete(apps[host], appName)
			LogVerbose("Removed app with name: " + appName + " and id: " + appID + " @" + host + " from known apps")
		}
	} else {
		if apps[host] == nil {
			apps[host] = map[string]string{}
		}
		apps[host][appName] = appID
		LogVerbose("Added app with name: " + appName + " and id: " + appID + " @" + host + " to known apps")
	}

	// Write to knownApps.yml
	out, _ := yaml.Marshal(apps)

	if err := ioutil.WriteFile(knownAppsFilePath, out, 0644); err != nil {
		FatalErrorf("could not write to '%s': %s", knownAppsFilePath, err)
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
				FatalError("could not create .corectl folder in home directory: ", err)
			}
		}

		// Create knownApps.yml in .corectl folder
		_, err := os.Create(knownAppsFilePath)
		if err != nil {
			FatalErrorf("could not create %s: %s", knownAppsFilePath, err)
		}

		LogVerbose("Created ~/.corectl/knownApps.yml for storage of app ids")
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
