package rest

import "C"
import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/qlik-oss/corectl/internal/log"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
)

type RestCallerSettings interface {
	TlsConfig() *tls.Config
	Insecure() bool
	Headers() http.Header
	RestBaseUrl() *neturl.URL
	AppId() string
	RestAdaptedAppId() string
	PrintMode() log.PrintMode
}

type RestCaller struct {
	RestCallerSettings
}

// parseFunction is any function that reads bytes into a pointer and returns an error,
// json.Unmarshal is an example.
type parseFunction func([]byte, interface{}) error

// Call performs the specified the request and uses the passed parsing function 'read'
// to parse the response into the supplied result.
// If the response status code is not 200 series an error will be returned.
// The supplied http request will be modified and should not be reused
func (c *RestCaller) Call(req *http.Request, result interface{}, read parseFunction) error {
	data, err := c.CallRaw(req)
	if err != nil {
		return err
	}
	err = read(data, result)
	if err != nil {
		return err
	}
	return nil
}

func (c *RestCaller) CallRaw(req *http.Request) ([]byte, error) {
	client := http.DefaultClient
	certs := c.TlsConfig()
	if certs != nil {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: certs,
			},
		}
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}
	for key := range c.Headers() {
		req.Header.Set(key, c.Headers().Get(key))
	}

	//t0 := time.Now()
	response, err := client.Do(req)
	//t1 := time.Now()
	//interval := t1.Sub(t0)
	//fmt.Println("Time", interval)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("got status %d", response.StatusCode)
	}
	data, _ := ioutil.ReadAll(response.Body)
	return data, nil
}

func (c *RestCaller) CallRest(method, path string, query map[string]string, body []byte) []byte {
	url := c.CreateUrl(path, query)
	method = strings.ToUpper(method)

	fmt.Fprintln(os.Stderr, method+": "+url.String())
	if len(body) > 0 {
		fmt.Fprintln(os.Stderr, string(body))
	}
	var req *http.Request
	var err error
	if len(body) == 0 {
		req, err = http.NewRequest(method, url.String(), nil)
	} else {
		buf := bytes.NewBuffer(body)
		req, err = http.NewRequest(method, url.String(), buf)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating request: ", err)
		os.Exit(1)
	}
	req.Header = c.Headers().Clone()
	b, err := c.CallRaw(req)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while sending request: ", err)
		os.Exit(1)
	}
	out := bytes.NewBuffer([]byte{}) // Indent because we are nice.
	if err := json.Indent(out, b, "", "  "); err == nil {
		b = out.Bytes()
	}
	// Make sure the bytes end with a newline, it is prettier.
	if l := len(b); l > 0 && b[l-1] != byte('\n') {
		b = append(b, byte('\n'))
	}
	return b
}

func (c *RestCaller) CreateUrl(path string, q map[string]string) *neturl.URL {

	url := c.RestBaseUrl()
	if strings.HasSuffix(url.Path, "/") {
		if strings.HasPrefix(path, "/") {
			url.Path += path[1:]
		} else {
			url.Path += path
		}
	} else {
		if strings.HasPrefix(path, "/") {
			url.Path += path
		} else {
			url.Path += "/" + path
		}
	}
	query := neturl.Values{}
	for k, v := range q {
		query.Add(k, v)
	}
	url.RawQuery = query.Encode()
	return url
}

func (c *RestCaller) CallGet(urlFormat string, q map[string]string, ids ...interface{}) ([]byte, error) {

	url := c.CreateUrl(fmt.Sprintf(urlFormat, ids...), q)
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	data, err := c.CallRaw(req)
	if err != nil {
		return nil, err
	}
	return data, err
}
