package internal

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var contextFilePath = path.Join(userHomeDir(), ".corectl", "contexts.yml")

func AddContext(contextName string, productName string, comment string) {
	createContextFileIfNotExist()
	_, contexts := GetContexts()

	if contextExists(contexts, contextName) {
		FatalErrorf("context with name '%s' already exists", contextName)
	}

	mymap := map[interface{}]interface{}{}
	mymap["engine"] = viper.GetString("engine")
	mymap["headers"] = viper.GetStringMapString("headers")
	mymap["certificates"] = viper.GetString("certificates")
	mymap["product"] = productName
	mymap["comment"] = comment

	contexts[contextName] = mymap

	LogVerbose("Added context with name: " + contextName)

	setContexts(contextName, contexts)
}

func RemoveContext(contextName string) {
	createContextFileIfNotExist()
	currentContext, contexts := GetContexts()

	if !contextExists(contexts, contextName) {
		FatalErrorf("context with name '%s' does not exist", contextName)
	}

	delete(contexts, contextName)
	LogVerbose("Removed context with name: " + contextName)

	if currentContext == contextName {
		setContexts("", contexts)
	} else {
		setContexts(currentContext, contexts)
	}
}

func contextExists(contexts map[interface{}]interface{}, contextName string) bool {
	if _, exists := contexts[contextName]; exists {
		LogVerbose("Found context: " + contextName)
		return true
	}
	return false
}

func SetCurrentContext(contextName string) {
	currentContext, contexts := GetContexts()

	if !contextExists(contexts, contextName) {
		FatalErrorf("context with name '%s' does not exist", contextName)
	}

	if currentContext == contextName {
		LogVerbose("Current context already set to " + contextName)
		return
	}

	LogVerbose("Set current context to: " + contextName)

	setContexts(contextName, contexts)
}

func setContexts(currentContext string, contexts map[interface{}]interface{}) {
	contexts["current-context"] = currentContext

	out, _ := yaml.Marshal(contexts)

	if err := ioutil.WriteFile(contextFilePath, out, 0644); err != nil {
		FatalErrorf("could not write to '%s': %s", contextFilePath, err)
	}
}

func GetContexts() (string, map[interface{}]interface{}) {
	var contexts = map[interface{}]interface{}{}
	yamlFile, err := ioutil.ReadFile(contextFilePath)
	if err != nil {
		return "", nil
	}
	err = yaml.Unmarshal(yamlFile, &contexts)
	if err != nil {
		FatalErrorf("could not parse content of contexts yaml '%s': %s", yamlFile, err)
	}

	if len(contexts) == 0 {
		return "", map[interface{}]interface{}{}
	}

	currentContext := contexts["current-context"].(string)
	delete(contexts, "current-context")

	return currentContext, contexts
}

func GetCurrentContext() map[interface{}]interface{} {
	currentContext, contexts := GetContexts()
	if currentContext == "" {
		return nil
	}
	res, _ := contexts[currentContext].(map[interface{}]interface{})
	return res
}

func GetSpecificContext(contextName string) map[interface{}]interface{} {
	_, contexts := GetContexts()

	if !contextExists(contexts, contextName) {
		return nil
	}

	res, _ := contexts[contextName].(map[interface{}]interface{})
	return res
}

// Create a contexts.yml if one does not exist
func createContextFileIfNotExist() {
	if _, err := os.Stat(contextFilePath); os.IsNotExist(err) {

		// Create .corectl folder in home directory
		if _, err := os.Stat(path.Join(userHomeDir(), ".corectl")); os.IsNotExist(err) {
			err = os.Mkdir(path.Join(userHomeDir(), ".corectl"), os.ModePerm)
			if err != nil {
				FatalError("could not create .corectl folder in home directory: ", err)
			}
		}

		// Create contexts.yml in .corectl folder
		_, err := os.Create(contextFilePath)
		if err != nil {
			FatalErrorf("could not create %s: %s", contextFilePath, err)
		}

		LogVerbose("Created ~/.corectl/contexts.yml for storage of corectl contexts")
	}
}
