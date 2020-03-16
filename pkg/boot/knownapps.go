package boot

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/qlik-oss/corectl/pkg/log"
	"gopkg.in/yaml.v2"
)

var knownAppsFilePath = path.Join(userHomeDir(), dynconf.ContextDir, "knownApps.yml")

// Fetch a matching app id from known apps for a specified app name
// If not found return the appName and found bool set to false
func ApplyNameToIDTransformation(host string, appName string) (appID string, found bool) {
	apps := GetKnownApps()

	if apps == nil {
		log.Verboseln("knownApps yaml file not found")
		return appName, false
	}

	if id, exists := apps[host][appName]; exists {
		log.Verboseln("Found id: " + id + " for app with name: " + appName + " @" + host)
		return id, true
	}

	log.Verboseln("No known id for app with name: " + appName)
	return appName, false
}

// Get map of known apps
func GetKnownApps() map[string]map[string]string {
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
func SetAppIDToKnownApps(host string, appName string, appID string, remove bool) {

	createKnownAppsFileIfNotExist()
	apps := GetKnownApps()

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

		// Create context folder in home directory
		corectlDir := path.Join(userHomeDir(), dynconf.ContextDir)
		if _, err := os.Stat(corectlDir); os.IsNotExist(err) {
			err = os.Mkdir(corectlDir, os.ModePerm)
			if err != nil {
				log.Fatalf("could not create %s folder in home directory: %s", dynconf.ContextDir, err)
			}
		}

		// Create knownApps.yml in context folder
		_, err := os.Create(knownAppsFilePath)
		if err != nil {
			log.Fatalf("could not create %s: %s\n", knownAppsFilePath, err)
		}

		log.Verbosef("Created ~/%s/knownApps.yml for storage of app ids", dynconf.ContextDir)
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
