package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/qlik-oss/corectl/internal/log"
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
// It also keeps a user's comments regarding the context.
type Context struct {
	Engine       string
	Headers      map[string]string
	Certificates string
	Comment      string
}

var contextFilePath = path.Join(userHomeDir(), ".corectl", "contexts.yml")

func SetContext(contextName, comment string) string {
	if contextName == "" {
		log.Fatalln("context name not supplied")
	}

	createContextFileIfNotExist()
	handler := NewContextHandler()

	var context *Context
	var update bool

	if handler.Exists(contextName) {
		context = handler.Get(contextName)
		log.Debugln("Updating context: " + contextName)
		update = true
	} else {
		context = &Context{}
		log.Debugln("Creating context: " + contextName)
	}

	updated := context.Update(&map[string]interface{}{
		"engine":       viper.GetString("engine"),
		"headers":      viper.GetStringMapString("headers"),
		"certificates": viper.GetString("certificates"),
		"comment":      comment,
	})

	if update {
		log.Debugf("Updated fields %v of context %s\n", updated, contextName)
	}

	if err := context.Validate(); err != nil {
		log.Fatalf("context '%s' is not valid: %s", contextName, err.Error())
	}

	if !update {
		handler.Contexts[contextName] = context
	}

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

func UseContext(contextName string) string {
	handler := NewContextHandler()
	handler.Use(contextName)
	handler.Save()
	return contextName
}

func ClearContext() string {
	handler := NewContextHandler()
	previous := handler.Clear()
	handler.Save()
	return previous
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
		log.Fatalf("could not parse content of contexts yaml '%s': %s", yamlFile, err)
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
		log.Debugln("Found context: " + contextName)
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

func (ch *ContextHandler) Use(contextName string) {
	if !ch.Exists(contextName) {
		log.Fatalf("context with name '%s' does not exist", contextName)
	}
	if ch.Current == contextName {
		log.Debugln("Current context already set to " + contextName)
		return
	}
	log.Debugln("Set current context to: " + contextName)

	ch.Current = contextName
}

func (ch *ContextHandler) Clear() (previous string) {
	if ch.Current == "" {
		log.Debugln("No context is set")
		return ""
	}
	previous = ch.Current
	log.Debugf("Unset current context '%s'\n", previous)
	ch.Current = ""
	return
}

func (ch *ContextHandler) Remove(contextName string) (string, bool) {
	if !ch.Exists(contextName) {
		log.Fatalf("context with name '%s' does not exist", contextName)
	}
	delete(ch.Contexts, contextName)
	log.Debugln("Removed context with name: " + contextName)
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
		log.Fatalf("could not write to '%s': %s", contextFilePath, err)
	}
}

// Update uses reflection to update a Context's fields.
// This method ignores empty strings and nil values so it will
// only update the context with new information provided.
// It returns the names of the updated fields.
func (c *Context) Update(m *map[string]interface{}) []string {
	ptr := reflect.ValueOf(c)
	val := reflect.Indirect(ptr)
	updated := []string{}
	for k, v := range *m {
		f := val.FieldByName(strings.Title(k))
		if f.IsValid() {
			vval := reflect.ValueOf(v)
			if hasValue(vval) && f.Type() == vval.Type() {
				f.Set(vval)
				updated = append(updated, k)
			}
		}
	}
	return updated
}

func (c *Context) ToMap() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	m["engine"] = c.Engine
	m["headers"] = c.Headers
	m["certificates"] = c.Certificates
	m["comment"] = c.Comment
	return m
}

func (c *Context) Validate() error {
	if c.Engine == "" && c.Certificates == "" && (c.Headers == nil || len(c.Headers) == 0) {
		return fmt.Errorf("empty context: no engine url, certificates path or headers specified, need at least one")
	}
	return nil
}

func hasValue(v reflect.Value) bool {
	switch k := v.Kind(); k {
	case reflect.String:
		return v.String() != ""
	case reflect.Map, reflect.Struct, reflect.Slice:
		return !v.IsNil()
	}
	return false
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
				log.Fatalln("could not create .corectl folder in home directory: ", err)
			}
		}

		// Create contexts.yml in .corectl folder
		_, err := os.Create(contextFilePath)
		if err != nil {
			log.Fatalf("could not create %s: %s", contextFilePath, err)
		}

		log.Debugln("Created ~/.corectl/contexts.yml for storage of corectl contexts")
	}
}
