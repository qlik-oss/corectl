package internal

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"syscall"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
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

// SetContext sets up the context to be used while communitcating with the engine
func SetContext(contextName, comment string) string {
	if contextName == "" {
		log.Fatalln("context name not supplied")
	}

	createContextFileIfNotExist()
	handler := NewContextHandler()

	var context *Context
	var update bool
	var certificates string

	if handler.Exists(contextName) {
		context = handler.Get(contextName)
		log.Verboseln("Updating context: " + contextName)
		update = true
	} else {
		context = &Context{}
		log.Verboseln("Creating context: " + contextName)
	}

	if certPath := viper.GetString("certificates"); certPath != "" {
		certificates = RelativeToProject(viper.GetString("certificates"))
	}

	updated := context.Update(&map[string]interface{}{
		"engine":       viper.GetString("engine"),
		"headers":      viper.GetStringMapString("headers"),
		"certificates": certificates,
		"comment":      comment,
	})

	if update {
		log.Verbosef("Updated fields %v of context %s\n", updated, contextName)
	}

	if err := context.Validate(); err != nil {
		log.Fatalf("context '%s' is not valid: %s\n", contextName, err.Error())
	}

	if !update {
		handler.Contexts[contextName] = context
	}

	handler.Current = contextName
	handler.Save()
	return contextName
}

// RemoveContext from context file
func RemoveContext(contextName string) (string, bool) {
	handler := NewContextHandler()
	contextName, wasCurrent := handler.Remove(contextName)
	handler.Save()
	return contextName, wasCurrent
}

// UseContext sets the current context based on name
func UseContext(contextName string) string {
	handler := NewContextHandler()
	handler.Use(contextName)
	handler.Save()
	return contextName
}

// ClearContext unsets the current context
func ClearContext() string {
	handler := NewContextHandler()
	previous := handler.Clear()
	handler.Save()
	return previous
}

// LoginContext login to a Qlik Sense Enterprise and sets the X-Qlik-Session as a cookie
func LoginContext(tlsClientConfig *tls.Config, contextName string) {
	userName := viper.GetString("user")
	password := viper.GetString("password")

	handler := NewContextHandler()
	var context *Context

	if contextName == "" {
		context = handler.GetCurrent()
		if context == nil {
			log.Fatalf(" no 'current-context' found in config.\n")
		}
		contextName = handler.Current
	} else {
		context = handler.Get(contextName)
		if context == nil {
			log.Fatalf("context '%s' wasn't found.\n", contextName)
		}
	}

	log.Infof("Using context '%s', with URL '%s'\n", contextName, context.Engine)

	qlikSession := getSessionCookie(tlsClientConfig, context.Engine, userName, password)

	if _, ok := context.Headers["cookie"]; ok {
		// Cookie header present
		re := regexp.MustCompile(`X-Qlik-Session=[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
		if re.MatchString(context.Headers["cookie"]) {
			context.Headers["cookie"] = re.ReplaceAllString(context.Headers["cookie"], qlikSession)
		} else {
			context.Headers["cookie"] = fmt.Sprintf("%s; %s", context.Headers["cookie"], qlikSession)
		}
	} else {
		// Cookie header has to be added
		if context.Headers == nil {
			context.Headers = map[string]string{}
		}
		context.Headers["cookie"] = qlikSession
	}

	handler.Save()
}

// NewContextHandler helps with handeling contexts
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
		log.Fatalf("could not parse content of contexts yaml '%s': %s\n", yamlFile, err)
	}

	if handler.Contexts == nil {
		handler.Contexts = map[string]*Context{}
	}

	if len(handler.Contexts) == 0 {
		return handler
	}
	return handler
}

// Exists checks if context exists
func (ch *ContextHandler) Exists(contextName string) bool {
	if _, ok := ch.Contexts[contextName]; ok {
		log.Verboseln("Found context: " + contextName)
		return ok
	}
	return false
}

// Get returns context if present
func (ch *ContextHandler) Get(contextName string) *Context {
	if context, ok := ch.Contexts[contextName]; ok {
		return context
	}
	return nil
}

// GetCurrent returns the context marked as current
func (ch *ContextHandler) GetCurrent() *Context {
	cur := ch.Current
	if cur == "" {
		return nil
	}
	return ch.Get(cur)
}

// Use sets the current context
func (ch *ContextHandler) Use(contextName string) {
	if !ch.Exists(contextName) {
		log.Fatalf("context with name '%s' does not exist\n", contextName)
	}
	if ch.Current == contextName {
		log.Verboseln("Current context already set to " + contextName)
		return
	}
	log.Verboseln("Set current context to: " + contextName)

	ch.Current = contextName
}

// Clear the context
func (ch *ContextHandler) Clear() (previous string) {
	if ch.Current == "" {
		log.Verboseln("No context is set")
		return ""
	}
	previous = ch.Current
	log.Verbosef("Unset current context '%s'\n", previous)
	ch.Current = ""
	return
}

// Remove the context
func (ch *ContextHandler) Remove(contextName string) (string, bool) {
	if !ch.Exists(contextName) {
		log.Fatalf("context with name '%s' does not exist\n", contextName)
	}
	delete(ch.Contexts, contextName)
	log.Verboseln("Removed context with name: " + contextName)
	wasCurrent := false
	if ch.Current == contextName {
		ch.Current = ""
		wasCurrent = true
	}
	return contextName, wasCurrent
}

// Save the context file
func (ch *ContextHandler) Save() {
	out, _ := yaml.Marshal(*ch)

	if err := ioutil.WriteFile(contextFilePath, out, 0644); err != nil {
		log.Fatalf("could not write to '%s': %s\n", contextFilePath, err)
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

// ToMap returns a map from the context
func (c *Context) ToMap() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	m["engine"] = c.Engine
	m["headers"] = c.Headers
	m["certificates"] = c.Certificates
	m["comment"] = c.Comment
	return m
}

// Validate that at least one property is set for the context
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
			log.Fatalf("could not create %s: %s\n", contextFilePath, err)
		}

		log.Verboseln("Created ~/.corectl/contexts.yml for storage of corectl contexts")
	}
}

func getSessionCookie(tlsClientConfig *tls.Config, engineURL string, userName string, password string) string {
	// Verify Qlik Sense URL
	u, err := url.Parse(engineURL)

	if err != nil {
		log.Fatalln("The engineURL doesn't seem to be correct")
	}

	if u.Scheme != "https" {
		log.Fatalln("Only login through secure connections (HTTPS) is supported")
	}

	// Get username
	if userName == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Username (domain\\user): ")
		userName, _ = reader.ReadString('\n')
	}

	if !strings.Contains(userName, "\\") {
		log.Fatalln("username MUST be in format 'domain\\user'")
	}

	// Get password
	if password == "" {
		fmt.Print("Enter Password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password = string(bytePassword)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsClientConfig

	// Get post URL via redirects
	resp, err := http.Get(engineURL)

	if err != nil {
		log.Fatalln(err)
	}

	// Generate xrfkey
	xrfkey := generateXrfkey()

	loginURL := resp.Request.URL
	q := loginURL.Query()
	q.Add("xrfkey", xrfkey)
	loginURL.RawQuery = q.Encode()

	urlData := url.Values{}
	urlData.Set("username", userName)
	urlData.Set("pwd", password)

	hc := http.Client{}
	req, err := http.NewRequest("POST", loginURL.String(), strings.NewReader(urlData.Encode()))

	if err != nil {
		log.Fatalln(err)
	}

	req.PostForm = urlData
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("x-qlik-xrfkey", xrfkey)

	postResp, err := hc.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	setCookie := postResp.Header.Get("Set-Cookie")

	if setCookie == "" {
		log.Fatalln("Not able to get the 'X-Qlik-Session' cookie, please check your password.")
	}

	return strings.TrimRight(strings.Fields(setCookie)[0], ";")
}

func generateXrfkey() string {

	b := make([]byte, 8)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	return fmt.Sprintf("%X", b)
}
