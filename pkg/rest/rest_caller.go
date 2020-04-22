package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/qlik-oss/corectl/pkg/log"
)

type loggableBody struct {
	io.ReadCloser
	content []byte
}

func (lb *loggableBody) String() string {
	return string(lb.content)
}

func (c *RestCaller) CreateLoggableJsonBody(data []byte) io.ReadCloser {
	buffer := ioutil.NopCloser(bytes.NewBuffer(data))
	res := loggableBody{
		ReadCloser: buffer,
		content:    data,
	}
	return res
}

type RestCallerSettings interface {
	TlsConfig() *tls.Config
	Insecure() bool
	Headers() http.Header
	RestBaseUrl() *neturl.URL
	AppId() string
	RestAdaptedAppId() string
	PrintMode() log.PrintMode
	IsSenseForKubernetes() bool
}

type RestCaller struct {
	RestCallerSettings
}

// Call builds and peforms the request defined by the parameters and parses the end result into the supplied result interface.
// If the response status code is not 200 series an error will be returned.
// The supplied http request may be modified and should not be reused.
func (c *RestCaller) CallStd(method, path, contentType string, queryParams map[string]string, body io.ReadCloser, result interface{}) error {
	url := c.CreateUrl(path, queryParams)
	req, err := http.NewRequest(strings.ToUpper(method), url.String(), body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	err = c.CallReq(req, result)
	return err
}

// Call performs the specified request and parses the end result into the supplied result interface
// If the response status code is not 200 series an error will be returned.
// The supplied http request may be modified and should not be reused
func (c *RestCaller) Call(method, path string, query map[string]string, body []byte) []byte {
	url := c.CreateUrl(path, query).String()
	var req *http.Request
	var err error
	if l := len(body); l == 0 {
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
	err = c.CallReq(req, &result)
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
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return BuildErrorWithBody(res, data)
	}
	if err != nil {
		return err
	}

	if result != nil {
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
	}
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
	log.Verbosef("%s %s", req.Method, req.URL)
	if c.PrintMode().VerboseMode() {
		logHeader(req.Header, "> ")
	}

	// TODO support logging more than only JSON?
	var buf *bytes.Buffer
	if req.Body != nil {
		if c.PrintMode().VerboseMode() {
			contentType := req.Header.Get("Content-Type")
			if contentType == "" || contentType == "application/json" {
				// Replace req.Body with a TeeReader which writes to buf on reads so we can log it.
				buf = bytes.NewBuffer([]byte{})
				req.Body = ioutil.NopCloser(io.TeeReader(req.Body, buf))
			}
		}
	}

	t0 := time.Now()
	response, err := client.Do(req)
	t1 := time.Now()

	if buf != nil {
		str := log.FormatAsJSON(buf.Bytes())
		if str != "" {
			log.Verbose("PAYLOAD:")
			log.Verbose(str)
		}
	}
	if c.PrintMode().VerboseMode() {
		logHeader(response.Header, "< ")
	}

	log.Verbose("Time ", t1.Sub(t0))
	if err != nil {
		return response, err
	}
	return response, nil
}

// logHeader logs a header (verbose) with the specified prefix.
func logHeader(header http.Header, prefix string) {
	keys := make([]string, len(header))
	i := 0
	for k := range header {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		if key == "Authorization" {
			value := header.Get(key)
			if strings.HasPrefix("Bearer ", value) {
				log.Verbosef("%s%s: %s", prefix, key, "Bearer **omitted**")
			} else {
				log.Verbosef("%s%s: %s", prefix, key, value)
			}
		} else {
			log.Verbosef("%s%s: %s", prefix, key, header.Get(key))
		}
	}
}
