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
func applyNameToIDTransformation(engineURL string, appName string) (appID string, found bool) {
	apps := getKnownApps()

	if apps == nil {
		LogVerbose("knownApps yaml file not found")
		return appName, false
	}

	if id, ok := apps[engineURL][appName]; ok {
		LogVerbose("Found id: " + id + " for app with name: " + appName + " @" + engineURL)
		return id, true
	}

	LogVerbose("No id found for app with name: " + appName)
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
		FatalError("Failed to unmarshal content of knownApps yaml file: " + err.Error())
	}

	return knownApps
}

// Add an app or remove an app from known apps
func setAppIDToKnownApps(engineURL string, appName string, appID string, remove bool) {

	createKnownAppsFileIfNotExist()
	apps := getKnownApps()

	// Either remove or add an entry
	if remove {
		_, exists := apps[engineURL][appName]
		if exists {
			delete(apps[engineURL], appName)
			LogVerbose("Removed app with name: " + appName + " and id: " + appID + " @" + engineURL + " from known apps")
		}
	} else {
		if apps[engineURL] == nil {
			apps[engineURL] = make(map[string]string)
		}
		apps[engineURL][appName] = appID
		LogVerbose("Added app with name: " + appName + " and id: " + appID + " @" + engineURL + " to known apps")
	}

	// Write to knownApps.yml
	out, _ := yaml.Marshal(apps)
	_ = ioutil.WriteFile(knownAppsFilePath, out, 0644)
}

// Create a knownApps.yml if one does not exist
func createKnownAppsFileIfNotExist() {
	if _, err := os.Stat(knownAppsFilePath); os.IsNotExist(err) {

		// Create .corectl folder in home directory
		err = os.Mkdir(path.Join(userHomeDir(), ".corectl"), os.ModePerm)
		if err != nil {
			FatalError("Failed to create .corectl folder in home directory: " + err.Error())
		}

		// Create knownApps.yml in .corectl folder
		_, err := os.Create(knownAppsFilePath)
		if err != nil {
			FatalError("Failed to create knownApps.yml in .corectl folder: " + err.Error())
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
