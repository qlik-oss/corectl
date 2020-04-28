package standard

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestGetJsonBodyFromFile(t *testing.T) {
	for i, testcase := range []struct {
		fileName         string
		bodyString       string
		params           map[string]string
		expectedBody     string
		expectedMimeType string
	}{
		{"raw_test_body.json", "", nil, "{\"test\": \"yes it is\"}\n", "application/json"},
		{"raw_test_body.qvf", "", nil, "<qvf content placeholder>\n", "binary/octet-stream"},
		{"raw_test_body.png", "", nil, "<skip assert>", "image/png"},
		{"", "inline-body-string", nil, "inline-body-string", "text/plain"},
		{"", "{\"test\": \"yes it is\"}", nil, "{\"test\": \"yes it is\"}", "application/json"},
		{"", "", map[string]string{"key": "value"}, "{\"key\":\"value\"}", "application/json"},
		{"", "", map[string]string{"key:bool": "true"}, "{\"key\":true}", "application/json"},
		{"", "", map[string]string{"nested.key": "value"}, "{\"nested\":{\"key\":\"value\"}}", "application/json"},
	} {
		body, mineType, err := getBodyFromFlags(testcase.fileName, testcase.bodyString, testcase.params)
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

func TestGetQvfBodyFromFile(t *testing.T) {
	body, mineType, err := getBodyFromFlags("raw_test_body.qvf", "", make(map[string]string))
	assert.Equal(t, "<qvf content placeholder>", body)
	assert.Equal(t, "binary/octet-stream", mineType)
	assert.Nil(t, err)
}

func TestGetImageBodyFromFile(t *testing.T) {
	body, mineType, err := getBodyFromFlags("raw_test_body.png", "", make(map[string]string))
	assert.Equal(t, "wef", body)
	assert.Equal(t, "image/png", mineType)
	assert.Nil(t, err)
}
