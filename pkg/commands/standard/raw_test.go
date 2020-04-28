package standard

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestGetJsonBodyFromFlags(t *testing.T) {
	for i, testcase := range []struct {
		fileName         string
		bodyString       string
		params           map[string]string
		expectedBody     string
		expectedMimeType string
	}{
		{"raw_test_body.json", "", nil, "{\"test\": \"yes it is\"}\n", "application/json"},
		{"raw_test_body.qvf", "", nil, "<qvf content placeholder>\n", "binary/octet-stream"},
		{"raw_test_body-image.png", "", nil, "<skip assert>", "image/png"},
		{"", "inline-body-string", nil, "inline-body-string", "text/plain"},
		{"", "{\"test\": \"yes it is\"}", nil, "{\"test\": \"yes it is\"}", "application/json"},
		{"", "", map[string]string{"key": "value"}, "{\"key\":\"value\"}", "application/json"},
		{"", "", map[string]string{"key:bool": "true"}, "{\"key\":true}", "application/json"},
		{"", "", map[string]string{"nested.key": "value"}, "{\"nested\":{\"key\":\"value\"}}", "application/json"},
	} {
		body, mineType, err := getBodyFromFlags(relativeToTestCase(testcase.fileName), testcase.bodyString, testcase.params)
		bodyBytes, _ := ioutil.ReadAll(body)
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if testcase.expectedBody != "<skip assert>" {
				assert.Equal(t, testcase.expectedBody, string(bodyBytes))
			}
			assert.Equal(t, testcase.expectedMimeType, mineType)
			assert.Nil(t, err)
		})
	}
}

func relativeToTestCase(fileName string) string {
	if fileName == "" {
		return ""
	}
	_, testCaseFilename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(testCaseFilename), fileName)
}
