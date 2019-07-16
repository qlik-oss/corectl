package internal

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ContextHandler struct {
	Current  string `yaml:"current-context"`
	Contexts map[string]*Context
}

type Context struct {
	Engine       string
	Headers      map[string]string
	Certificates string
	Product      string
	Comment      string
}

var contextFilePath = path.Join(userHomeDir(), ".corectl", "contexts.yml")

func AddContext(contextName string, productName string, comment string) {
	if contextName == "" {
		FatalError("\"\" is not a valid context name")
	}
	createContextFileIfNotExist()
	handler := NewContextHandler()

	if handler.Exists(contextName) {
		FatalErrorf("context with name '%s' already exists", contextName)
	}

	context := &Context{
		Engine:       viper.GetString("engine"),
		Headers:      viper.GetStringMapString("headers"),
		Certificates: viper.GetString("certificates"),
		Product:      productName,
		Comment:      comment,
	}

	handler.Contexts[contextName] = context

	LogVerbose("Added context with name: " + contextName)

	handler.Current = contextName
	handler.Save()
}

func RemoveContext(contextName string) {
	handler := NewContextHandler()
	handler.Remove(contextName)
	LogVerbose("Removed context with name: " + contextName)
	handler.Save()
}

func SetCurrentContext(contextName string) {
	handler := NewContextHandler()
	handler.SetCurrent(contextName)
	handler.Save()
}

func NewContextHandler() *ContextHandler {
	handler := &ContextHandler{}
	if !fileExists(contextFilePath) {
		return handler
	}
	yamlFile, err := ioutil.ReadFile(contextFilePath)
	if err != nil {
		return handler
	}
	err = yaml.Unmarshal(yamlFile, &handler)
	if err != nil {
		FatalErrorf("could not parse content of contexts yaml '%s': %s", yamlFile, err)
	}

	if handler.Contexts == nil {
		handler.Contexts = map[string]*Context{}
	}

	if len(handler.Contexts) == 0 {
		return handler
	}
	return handler
}

func (ch *ContextHandler) Exists(contextName string) bool {
	if _, ok := ch.Contexts[contextName]; ok {
		LogVerbose("Found context: " + contextName)
		return ok
	}
	return false
}

func (ch *ContextHandler) Get(contextName string) *Context {
	if context, ok := ch.Contexts[contextName]; ok {
		return context
	}
	return nil
}

func (ch *ContextHandler) GetCurrent() *Context {
	cur := ch.Current
	if cur == "" {
		return nil
	}
	return ch.Get(cur)
}

func (ch *ContextHandler) SetCurrent(contextName string) {
	if !ch.Exists(contextName) {
		FatalErrorf("context with name '%s' does not exist", contextName)
	}
	if ch.Current == contextName {
		LogVerbose("Current context already set to " + contextName)
		return
	}
	LogVerbose("Set current context to: " + contextName)

	ch.Current = contextName
}

func (ch *ContextHandler) Remove(contextName string) {
	if !ch.Exists(contextName) {
		FatalErrorf("context with name '%s' does not exist", contextName)
	}
	delete(ch.Contexts, contextName)
	LogVerbose("Removed context with name: " + contextName)
	if ch.Current == contextName {
		ch.Current = ""
	}
}

func (ch *ContextHandler) Save() {
	out, _ := yaml.Marshal(*ch)

	if err := ioutil.WriteFile(contextFilePath, out, 0644); err != nil {
		FatalErrorf("could not write to '%s': %s", contextFilePath, err)
	}
}

func (c *Context) ToMap() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	m["engine"] = c.Engine
	m["headers"] = c.Headers
	m["certificates"] = c.Certificates
	m["product"] = c.Product
	m["comment"] = c.Comment
	return m
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Create a contexts.yml if one does not exist
func createContextFileIfNotExist() {
	if !fileExists(contextFilePath) {

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
