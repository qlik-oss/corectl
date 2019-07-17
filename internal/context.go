package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ContextHandler maps strings to contexts and keeps track of the current context.
// It has various methods for manipulating and accessing Contexts.
type ContextHandler struct {
	Current  string `yaml:"current-context"`
	Contexts map[string]*Context
}

// Context represents a context. As of now it only contains information regarding connections.
// Meaning: engine url, certificates path and any headers.
// It also keeps the product it is meant to be used for as well as the user's comments
// regarding the context.
type Context struct {
	Engine       string
	Headers      map[string]string
	Certificates string
	Product      string
	Comment      string
}

// products contains a mapping from shorthand to proper name of Qlik products.
// This map should not be modified.
var products = map[string]string{
	"QC": "Qlik Core",
	"QSC": "Qlik Sense Cloud",
	"QSD": "Qlik Sense Desktop",
	"QSE": "Qlik Sense Enterprise",
	"QSEoK": "Qlik Sense Enterpries on Kubernetes",
	"QSEoW": "Qlik Sense Enterpries on Windows",
}

func sortKeys(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func GetProducts() string {
	keys := sortKeys(products)
	l := len(keys)
	str := make([]string, l)
	for i, k := range keys {
		str[i] = fmt.Sprintf("%s (%s)", k, products[k])
	}
	prods := strings.Join(str[:l - 1], ", ")
	prods += " or " + str[l - 1]
	return prods
}

func isProduct(p string) bool {
	if _, ok := products[p]; ok {
		return true
	}
	return false
}

var contextFilePath = path.Join(userHomeDir(), ".corectl", "contexts.yml")

func CreateContext(contextName string, productName string, comment string) string {
	if contextName == "" {
		FatalError("\"\" is not a valid context name")
	}
	if !isProduct(productName) {
		// How we print string arrays should be handled by some utils or such
		keys := sortKeys(products)
		for i, k := range keys {
			keys[i] = "\"" + k + "\""
		}
		FatalErrorf("no product by the name '%s', should be one of: [%s]", productName, strings.Join(keys, ", "))
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

	if err := context.Validate(); err != nil {
		FatalErrorf("cannot create context '%s': %s", contextName, err.Error())
	}

	handler.Contexts[contextName] = context

	LogVerbose("Added context with name: " + contextName)

	handler.Current = contextName
	handler.Save()
	return contextName
}

func RemoveContext(contextName string) (string, bool) {
	handler := NewContextHandler()
	contextName, wasCurrent := handler.Remove(contextName)
	handler.Save()
	return contextName, wasCurrent
}

func SetCurrentContext(contextName string) string {
	handler := NewContextHandler()
	handler.SetCurrent(contextName)
	handler.Save()
	return contextName
}

func UnsetCurrentContext() {
	handler := NewContextHandler()
	handler.UnsetCurrent()
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

func (ch *ContextHandler) UnsetCurrent() (previous string) {
	if ch.Current == "" {
		LogVerbose("No context is set")
		return ""
	}
	previous = ch.Current
	LogVerbose(fmt.Sprintf("Unset context '%s'", previous))
	ch.Current = ""
	return
}

func (ch *ContextHandler) Remove(contextName string) (string, bool) {
	if !ch.Exists(contextName) {
		FatalErrorf("context with name '%s' does not exist", contextName)
	}
	delete(ch.Contexts, contextName)
	LogVerbose("Removed context with name: " + contextName)
	wasCurrent := false
	if ch.Current == contextName {
		ch.Current = ""
		wasCurrent = true
	}
	return contextName, wasCurrent
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

func (c *Context) Validate() error {
	if c.Engine == "" && c.Certificates == "" && (c.Headers == nil || len(c.Headers) == 0) {
		return fmt.Errorf("empty context: no engine url, certificates path or headers specified, need at least one")
	}
	return nil
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
