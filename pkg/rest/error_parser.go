package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//restError is the errors returned from rest calls. It contains a human readable message plus the raw body
type (
	Error interface {
		error
		Body() []byte
		StatusCode() int
	}

	restError struct {
		status  int
		message string
		body    []byte
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
	msg := fmt.Sprintf("%d %s", e.status, http.StatusText(e.status))
	if e.message != "" {
		msg += ": " + e.message
	}
	return msg
}

func (e *restError) Body() []byte {
	return e.body
}

func (e *restError) StatusCode() int {
	return e.status
}

// NewError creates an error that has both a readable message and the original body
func NewError(res *http.Response) Error {
	var message string
	data, err := ioutil.ReadAll(res.Body)
	if err == nil {
		message = buildStandardErrorMessage(data)
		if message == "" {
			message = buildNonArrayErrorMessage(data)
		}
		if message == "" {
			message = buildOtherKnownJsonErrorMessages(data)
		}
		if message == "" {
			message = buildNonJsonMessage(data)
		}
	}
	restErr := &restError{
		status:  res.StatusCode,
		body:    data,
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
