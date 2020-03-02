package rest

import "C"
import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/qlik-oss/corectl/internal/log"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"
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

// Call performs the specified request and parses the end result into the supplied result interface
// If the response status code is not 200 series an error will be returned.
// The supplied http request may be modified and should not be reused
func (c *RestCaller) Call(method, path string, query map[string]string, body []byte) []byte {
	url := c.CreateUrl(path, query).String()
	var req *http.Request
	var err error
	if len(body) == 0 {
		req, err = http.NewRequest(method, url, nil)
	} else {
		buf := bytes.NewBuffer(body)
		req, err = http.NewRequest(method, url, buf)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating request: ", err)
		os.Exit(1)
	}

	var result []byte
	err = c.CallReq(req, result)
	return result
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

// CallReq performs the specified request and parses the end result into the supplied result interface
// If the response status code is not 200 series an error will be returned.
// The supplied http request may be modified and should not be reused
func (c *RestCaller) CallReq(req *http.Request, result interface{}) error {
	res, err := c.CallRaw(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	switch x := result.(type) {
	case *[]byte:
		*x = data
	case *json.RawMessage:
		*x = data
	case *string:
		*x = string(data)
	default:
		err = json.Unmarshal(data, result)
	}

	return err
}

// Call builds and peforms the request defined by the parameters and parses the end result into the supplied result interface.
// If the response status code is not 200 series an error will be returned.
// The supplied http request may be modified and should not be reused.
func (c *RestCaller) CallStd(method, path string, queryParams map[string]string, body io.ReadCloser, result interface{}) error {
	url := c.CreateUrl(path, queryParams)
	req, err := http.NewRequest(strings.ToUpper(method), url.String(), body)
	err = c.CallReq(req, result)
	return err
}

// Call performs the request and returns the response.
// Note that the body of the response must be closed.
func (c *RestCaller) CallRaw(req *http.Request) (*http.Response, error) {
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

	var t0 time.Time
	if log.Traffic {
		fmt.Fprintln(os.Stderr, req.Method+": "+req.URL.String())
		t0 = time.Now()
	}

	response, err := client.Do(req)
	if log.Traffic {
		t1 := time.Now()
		interval := t1.Sub(t0)
		fmt.Fprintln(os.Stderr, "Time", interval)
	}
	if err != nil {
		return response, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("got status %d", response.StatusCode)
	}
	return response, nil
}
