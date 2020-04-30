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
	location := res.Header.Get("Location")

	if location != "" && (res.StatusCode == 201 || 301 <= res.StatusCode && res.StatusCode <= 307) {
		redirectUrl := c.RestBaseUrl()
		redirectUrl.Path = location
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
func (c *RestCaller) CallStreaming(method string, path string, query map[string]string, mimeType string, body io.ReadCloser, output io.Writer, raw bool, quiet bool) error {
	// Create the request
	url := c.CreateUrl(path, query)
	req, err := http.NewRequest(strings.ToUpper(method), url.String(), body)
	req.Header.Set("Content-Type", mimeType)

	//Make the actual invocation
	res, err := c.CallRawAndFollowRedirect(req)
	if err != nil {
		fmt.Println(output, err)
		return err
	}
	defer res.Body.Close()

	//Something when wrong it seems
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintln(output, err)
			return err
		}
		errorWithBody := BuildErrorWithBody(res, data)
		if raw {
			fmt.Fprintln(output, string(errorWithBody.Body()))
		} else {
			fmt.Fprintln(output, errorWithBody.Error())
		}
		return err
	}

	switch contentType := getContentType(res); contentType {
	case "application/json":
		//We have got a json response
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintln(output, err)
			return err
		}
		if quiet { //Only print IDs
			fmt.Fprint(output, string(filterIdsOnly(data)))
		} else if !raw { //Print data payload neatly formatted
			fmt.Fprint(output, log.FormatAsJSON(filterOutputForPrint(data)))
		} else { // Print it all
			fmt.Fprint(output, string(data))
		}
	case "text/plain", "text/html":
		// We've got some sort of text, safe to stream it to the output
		io.Copy(output, res.Body)
	default:
		//We have got something else which we should probably not print to terminal
		if file, ok := output.(*os.File); ok {
			fi, _ := file.Stat()
			if fi.Mode()&os.ModeCharDevice == os.ModeCharDevice {
				fmt.Fprintf(output, `Error: Content of type %q cannot be written directly to the terminal
       as it may contain special characters which can mess with your terminal.
       Specify an output location instead, either by flag or by piping the output to a file.
`, contentType)
				return fmt.Errorf("%q cannot be written directly to the terminal", contentType)
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
			if strings.HasPrefix(value, "Bearer") {
				log.Verbosef("%s%s: %s", prefix, key, "Bearer **omitted**")
			} else {
				log.Verbosef("%s%s: %s", prefix, key, value)
			}
		} else {
			log.Verbosef("%s%s: %s", prefix, key, header.Get(key))
		}
	}
}

func getContentType(res *http.Response) string {
	if typ := res.Header.Get("Content-Type"); typ != "" {
		typ = strings.Split(typ, ";")[0] // Strip charset utf-8 and such meta-data
		return typ
	}
	p := make([]byte, 512)
	n, err := res.Body.Read(p)
	if err != nil {
		return "application/octet-stream"
	}
	p = p[:n]
	typ := http.DetectContentType(p)
	log.Verbosef("Detected Content-Type: %q", typ)
	buf := ioutil.NopCloser(bytes.NewBuffer(p))
	// Concatenate the body back together with a MultiReadCloser (as we want to be able to close
	// the body).
	res.Body = MultiReadCloser(buf, res.Body)
	return typ
}

type multiReadCloser struct {
	io.Reader
	closers []io.Closer
}

// MultiReadCloser creates an io.ReadCloser which works as an io.MultiReader with
// the addition of being able to close all contained io.ReadClosers.
func MultiReadCloser(readClosers ...io.ReadCloser) io.ReadCloser {
	readers := make([]io.Reader, len(readClosers))
	for i, readCloser := range readClosers {
		readers[i] = readCloser
	}
	closers := make([]io.Closer, len(readClosers))
	for i, readCloser := range readClosers {
		closers[i] = readCloser
	}
	return &multiReadCloser{io.MultiReader(readers...), closers}
}

func (r *multiReadCloser) Close() error {
	var errors []string
	for i, closer := range r.closers {
		if err := closer.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("(%d %T): %s", i, closer, err.Error()))
		}
	}
	if len(errors) != 0 {
		return fmt.Errorf("failed to close: %s", strings.Join(errors, ", "))
	}
	return nil
}

func isJsonResponse(res *http.Response) bool {
	contentType := res.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "application/json")
}

func filterIdsOnly(bytes []byte) []byte {
	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		log.Warn(err)
		return nil
	}
	var ids string
	if data, ok := result["data"].([]interface{}); ok {
		for _, obj := range data {
			if m, ok := obj.(map[string]interface{}); ok {
				ids += fmt.Sprint(m["id"]) + "\n"
			}
		}
	} else if id, ok := result["id"].(string); ok {
		ids += id + "\n"
	}
	return []byte(ids)
}

func filterOutputForPrint(bytes []byte) []byte {
	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return bytes
	}
	data := result["data"]
	if data != nil {
		if dataArray, ok := data.([]interface{}); ok {
			for _, item := range dataArray {
				removeLinks(item)
			}
		} else if dataMap, ok := data.(map[string]interface{}); ok {
			removeLinks(dataMap)
		}
		return marshal(data)
	}
	removeLinks(result)
	return marshal(result)
}

func marshal(tree interface{}) []byte {
	return []byte(log.FormatAsJSON(tree))
}

func removeLinks(resultRaw interface{}) {
	if result, ok := resultRaw.(map[string]interface{}); ok {
		if links, ok := result["links"].(map[string]interface{}); ok {
			if links["self"] != nil || links["next"] != nil || links["prev"] != nil || links["Self"] != nil || links["Next"] != nil || links["Prev"] != nil {
				delete(result, "links")
			}
		}
	}
}
