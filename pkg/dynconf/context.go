package dynconf

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
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/qlik-oss/corectl/pkg/log"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// ContextHandler maps strings to contexts and keeps track of the current context.
// It has various methods for manipulating and accessing Contexts.
type ContextHandler struct {
	Current  string `yaml:"current-context"`
	Contexts map[string]Context
}

// Context is a map[string]interface{} to allow us to put any information in it we
// deem useful. This will in most cases consist of connection information such as
// server URL, headers for authentication and some info such as name/description/comment
// for the user.
type Context map[string]interface{}

const ContextDir = ".qlik"

var contextFilePath = path.Join(userHomeDir(), ContextDir, "contexts.yml")

// CreateContext creates a new context with the specified name and data.
func CreateContext(contextName string, data map[string]interface{}) {

	if contextName == "" {
		log.Fatalln("context name not supplied")
	}

	createContextFileIfNotExist()
	handler := NewContextHandler()

	if handler.Exists(contextName) {
		log.Fatalf("Context '%s' already exists", contextName)
	}

	context := Context{}
	log.Verboseln("Creating context: " + contextName)
	updated := context.Update(&data)
	log.Verbosef("Set fields %v for context %s", updated, contextName)

	if err := context.Validate(); err != nil {
		log.Fatalf("context '%s' is not valid: %s\n", contextName, err.Error())
	}

	handler.Contexts[contextName] = context
	handler.Save()
}

// UpdateContext updates the specified context with the provided data.
func UpdateContext(contextName string, data map[string]interface{}) {
	if contextName == "" {
		log.Fatalln("context name not supplied")
	}

	createContextFileIfNotExist()
	handler := NewContextHandler()

	if !handler.Exists(contextName) {
		log.Fatalf("No context by the name '%s'", contextName)
	}

	context := handler.Get(contextName)
	log.Verboseln("Updating context: " + contextName)

	updated := context.Update(&data)

	log.Verbosef("Updated fields %v of context %s\n", updated, contextName)

	if err := context.Validate(); err != nil {
		log.Fatalf("context '%s' is not valid: %s\n", contextName, err.Error())
	}

	handler.Save()
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
func LoginContext(tlsClientConfig *tls.Config, contextName string, userName string, password string) {

	handler := NewContextHandler()
	var context Context

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

	server := context.GetString("server")
	if server == "" {
		log.Fatalf("Context '%s' does not have any URL specified", contextName)
	}

	log.Infof("Using context '%s', with URL '%s'\n", contextName, server)
	qlikSession := getSessionCookie(tlsClientConfig, server, userName, password)
	headers := context.Headers()
	if headers == nil {
		headers = map[string]string{}
		context["headers"] = headers
	}
	if cookie, ok := headers["cookie"]; ok {
		// Cookie header present
		re := regexp.MustCompile(`X-Qlik-Session=[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
		if re.MatchString(cookie) {
			headers["cookie"] = re.ReplaceAllString(cookie, qlikSession)
		} else {
			headers["cookie"] = fmt.Sprintf("%s; %s", cookie, qlikSession)
		}
	} else { // Cookie header has to be added
		headers["cookie"] = qlikSession
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
		handler.Contexts = map[string]Context{}
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
func (ch *ContextHandler) Get(contextName string) Context {
	if context, ok := ch.Contexts[contextName]; ok {
		return context
	}
	return nil
}

// GetCurrent returns the context marked as current
func (ch *ContextHandler) GetCurrent() Context {
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

// Update updates a Context's fields.
// This method ignores empty strings and nil values so it will
// only update the context with new information provided.
// It returns the names of the updated fields.
func (c Context) Update(m *map[string]interface{}) []string {
	updated := []string{}
	for k, v := range *m {
		if k != "" && v != nil { // Might need reflection for the interface{}
			c[k] = v
			updated = append(updated, k)
		}
	}
	return updated
}

// Validate that at least one property is set for the context
func (c Context) Validate() error {
	if h, ok := c["headers"]; ok {
		if x, ok := h.(map[string]string); !ok {
			log.Fatalf("%T: %v", x, x)
			return fmt.Errorf(`headers must be a map, e.g. "Authorization": "Bearer MyJWT"`)
		}
	}
	return nil
}

// Headers retrieves the headers as a type that is easier to handle than
// an interface{} => interface{} map.
func (c Context) Headers() map[string]string {
	if h, ok := c["headers"]; ok {
		if x, ok := h.(map[interface{}]interface{}); !ok {
			log.Fatalln("context field 'headers' was not a map")
		} else {
			// Have to convert interface{} => interface{} to string => string somehow.
			// ¯\_(ツ)_/¯
			headers := map[string]string{}
			for k, v := range x {
				str := strings.ToLower(k.(string))
				headers[str] = v.(string)
			}
			return headers
		}
	}
	return nil
}

func (c Context) GetString(key string) string {
	if v, ok := c[key]; ok {
		if val, ok := v.(string); !ok {
			log.Fatalf("context field '%s' was not a string", key)
		} else {
			return val
		}
	}
	return ""
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

		// Create context folder in home directory
		if _, err := os.Stat(path.Join(userHomeDir(), ContextDir)); os.IsNotExist(err) {
			err = os.Mkdir(path.Join(userHomeDir(), ContextDir), os.ModePerm)
			if err != nil {
				log.Fatalf("could not create %s folder in home directory: %s", ContextDir, err.Error())
			}
		}

		// Create contexts.yml in context folder
		_, err := os.Create(contextFilePath)
		if err != nil {
			log.Fatalf("could not create %s: %s\n", contextFilePath, err)
		}

		log.Verbosef("Created ~/%s/contexts.yml for storage of corectl contexts", ContextDir)
	}
}

func getSessionCookie(tlsClientConfig *tls.Config, serverURL string, userName string, password string) string {
	// Verify Qlik Sense URL
	u, err := url.Parse(serverURL)

	if err != nil {
		log.Fatalln("The serverURL doesn't seem to be correct: ", err)
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
	resp, err := http.Get(serverURL)

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
