package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"

	"github.com/qlik-oss/corectl/pkg/log"
)

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

func logJSONPayload(buf *bytes.Buffer) {
	str := log.FormatAsJSON(buf.Bytes())
	if str == "" {
		str = buf.String()
	}
	if str != "" {
		log.Verbose("PAYLOAD:")
		log.Verbose(str)
	}
}

func logMultipartPayload(buf *bytes.Buffer, boundary string) {
	log.Verbose("PAYLOAD (MULTIPART):")
	r := multipart.NewReader(buf, boundary)
	var part *multipart.Part
	var err error
	for err == nil {
		part, err = r.NextPart()
		if err != nil {
			break
		}
		b, err := ioutil.ReadAll(part)
		if err != nil {
			break
		}
		log.Verbose("------")
		k := "Content-Disposition"
		if v := part.Header.Get(k); v != "" {
			log.Verbosef("> %s: %s", k, v)
		}
		k = "Content-Type"
		var contentType string
		if v := part.Header.Get(k); v != "" {
			contentType = v
			log.Verbosef("> %s: %s", k, contentType)
		}
		var str string
		switch contentType {
		case "", "application/json":
			str = log.FormatAsJSON(b)
			if str == "" {
				str = string(b)
			}
		case "application/octet-stream":
			str = fmt.Sprintf("Binary data: %d bytes", len(b))
		default:
			str = string(b)
		}
		log.Verbose("PART:")
		log.Verbose(str)
	}
	log.Verbose("------")
	if err != io.EOF {
		log.Error("malformed multipart request: ", err)
	}
}

// getContentType returns the content type of a response.
// If the Content-Type field in the header is not set the function wil
// use http.DetectContentType on the first 512 bytes (max) to determine
// the content-type. The response body will not be consumed, but it's pointer
// will change.
//
// If the body is empty the content-type won't be changed, meaning it will be
// either empty ("") or whatever it was set to in the header.
func getContentType(res *http.Response) string {
	if typ := res.Header.Get("Content-Type"); typ != "" {
		typ = strings.Split(typ, ";")[0] // Strip charset utf-8 and such meta-data
		return typ
	}
	p := make([]byte, 512)
	n, err := res.Body.Read(p)

	// Empty body, return empty content-type "".
	if err == io.EOF || n == 0 {
		log.Verbosef("Empty response body")
		return ""
	}
	// Some non-EOF error occured, return "application/octet-stream" which is default.
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

// multiReadCloser wraps a MultiReader and implements its own Close function
// which calls the Close function of its contained io.ReadClosers.
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

// Close closes all underlying io.ReadClosers. It does not throw the first error
// encountered but a concatenation of all errors (or nil).
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

type Filter func([]byte) []byte

// QuietFilter extracts all "id" fields and prints them.
func QuietFilter(bytes []byte) []byte {
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

// RawFilter just returns its input.
func RawFilter(bytes []byte) []byte {
	return []byte(log.FormatAsJSON(bytes))
}

// StandardFilter removes information not deemed of interest in a CLI context,
// such as links.
func StandardFilter(bytes []byte) []byte {
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
