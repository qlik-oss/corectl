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
	"strconv"
	"strings"
	"time"

	"github.com/qlik-oss/corectl/pkg/log"
)

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

// CallRawAndFollowRedirect is a shorthand for wrapping the response of CallRaw within a call to FollowRedirect
func (c *RestCaller) CallRawAndFollowRedirect(req *http.Request) (*http.Response, error) {
	return c.FollowRedirect(c.CallRaw(req))
}

// FollowRedirect takes the output of a previous call to CallReq and makes another request IF the first one
// contains a location header and has a statusCode of 201 or 301-307. The server of the RestBaseUrl is used
// to build the actual url
func (c *RestCaller) FollowRedirect(res *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	length, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	location := res.Header.Get("Location")

	// Don't follow redirect if we have a payload and the response is 201 Created.
	// 301-307 responses MAY contain payloads as well but, as such responses are mainly
	// for redirects the payload is probably of little significance.
	if length != 0 && res.StatusCode == 201 {
		return res, nil
	}

	if location != "" && (res.StatusCode == 201 || 301 <= res.StatusCode && res.StatusCode <= 307) {
		// Location might not be relative but absolute.
		redirectUrl := c.RestBaseUrl()
		locationUrl, _ := neturl.Parse(location)
		if locationUrl.Scheme == "" {
			redirectUrl.Path = locationUrl.Path
		} else {
			redirectUrl = locationUrl
		}
		if res.Body != nil {
			res.Body.Close() //Close body of original response
		}
		req, err := http.NewRequest("GET", redirectUrl.String(), nil)
		if err != nil {
			return nil, err
		}
		return c.CallRaw(req)
	}
	return res, nil
}

//CallStreaming makes a rest call, follows redirects if present and streams the output to the supplied output
func (c *RestCaller) CallStreaming(method string, path string, query map[string]string, mimeType string, body io.ReadCloser, output io.Writer, filter Filter) error {
	// Create the request
	url := c.CreateUrl(path, query)
	req, err := http.NewRequest(strings.ToUpper(method), url.String(), body)

	// Don't set content-type if the passed MIMEtype is invalid.
	// An empty body does not need to have content-type set.
	if mimeType != "" {
		req.Header.Set("Content-Type", mimeType)
	}

	//Make the actual invocation
	res, err := c.CallRawAndFollowRedirect(req)
	if err != nil {
		fmt.Fprintln(output, err)
		return err
	}
	defer res.Body.Close()

	//Something when wrong it seems
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		errorWithBody := BuildErrorWithBody(res, data)
		return errorWithBody
	}

	log.Verbose(res.Status)

	switch contentType := getContentType(res); contentType {
	case "application/json":
		//We have got a json response
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if filter == nil {
			filter = StandardFilter
		}
		fmt.Fprint(output, string(filter(data)))
	case "", "text/plain", "text/html":
		// We've got some sort of text or an empty body, safe to write it to the output.
		io.Copy(output, res.Body)
	default:
		// We have got something else which we should probably not print to terminal.
		if file, ok := output.(*os.File); ok {
			fi, _ := file.Stat()
			// If the file mode contains ModeCharDevice it means it's a terminal.
			if fi.Mode()&os.ModeCharDevice != 0 {
				err := fmt.Errorf(`Error: Content of type %q cannot be written directly to the terminal
       as it may contain special characters which can mess with your terminal.
       Specify an output location instead, either by flag or by piping the output to a file.
`, contentType)
				return err
			}
		}
		io.Copy(output, res.Body)
	}
	return nil
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
			if contentType == "" || contentType == "application/json" || strings.Contains(contentType, "multipart/form-data") {
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
		if str == "" {
			str = buf.String()
		}
		if str != "" {
			log.Verbose("PAYLOAD:")
			log.Verbose(str)
		}
	}
	if c.PrintMode().VerboseMode() {
		logHeader(response.Header, "< ")
	}

	log.Verbosef("Response time: %dms", t1.Sub(t0).Milliseconds())
	if err != nil {
		return response, err
	}
	return response, nil
}
