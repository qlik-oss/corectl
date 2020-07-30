package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/qlik-oss/corectl/pkg/log"
)

//restError is the errors returned from rest calls. It contains a human readable message plus the raw body
type (
	Error interface {
		error
		StatusCode() int
	}

	restError struct {
		status  int
		message string
	}

	StandardErrorItem struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
		Title  string `json:"title"`
	}

	StandardErrorArray struct {
		Errors []StandardErrorItem `json:"errors"`
	}

	OtherErrorFormats struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
)

func (e *restError) Error() string {
	if e.message == "" {
		return "empty response"
	}
	return e.message
}

func (e *restError) StatusCode() int {
	return e.status
}

// NewError creates an error that has both a readable message and the original body
func NewError(res *http.Response) Error {
	var message string
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		message = "failed to read response body: " + err.Error()
	} else {
		contentType := res.Header.Get("Content-Type")
		contentType = strings.Split(contentType, "; ")[0]
		switch contentType {
		case "", "application/json":
			message = string(log.FormatAsJSON(body))
		case "text/html":
			fallthrough
		default:
			message = string(body)
		}
	}
	restErr := &restError{
		status:  res.StatusCode,
		message: message,
	}
	return restErr
}

func buildStandardErrorMessage(data []byte) string {
	var result StandardErrorArray
	err := json.Unmarshal(data, &result)
	if err != nil {
		return ""
	}
	if len(result.Errors) == 0 {
		return ""
	}
	message := ""
	for i, errItem := range result.Errors {
		if i > 0 {
			message += "; "
		}
		message += errorItemToString(errItem)
	}
	return message
}

func buildNonArrayErrorMessage(data []byte) string {
	var result StandardErrorItem
	err := json.Unmarshal(data, &result)
	if err != nil {
		return ""
	}
	return errorItemToString(result)
}

func errorItemToString(result StandardErrorItem) string {
	message := ""
	appendSubSectionToMessage(&message, result.Title)
	if result.Detail != result.Title {
		appendSubSectionToMessage(&message, result.Detail)
	}
	appendCode(&message, result.Code)
	return message
}

func buildOtherKnownJsonErrorMessages(data []byte) string {
	var result OtherErrorFormats
	err := json.Unmarshal(data, &result)
	if err != nil {
		return ""
	}
	message := ""
	appendSubSectionToMessage(&message, result.Error)
	appendSubSectionToMessage(&message, result.Message)
	appendCode(&message, result.Code)
	return message
}

func buildNonJsonMessage(data []byte) string {
	return string(data)
}

func appendSubSectionToMessage(message *string, section string) {
	if section != "" {
		if *message != "" {
			*message += ": "
		}
		*message += section
	}
}

func appendCode(message *string, code string) {
	if code != "" {
		if *message != "" {
			*message += " "
		}
		*message += "(" + code + ")"
	}
}
