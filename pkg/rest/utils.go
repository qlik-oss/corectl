package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
