package standard

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/pkg/rest"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const maxJsonFileSize = 1024 * 1024
const applicationJson = "application/json"
const binaryOctetStream = "binary/octet-stream"
const applicationOctetStream = "application/octet-stream"
const textPlain = "text/plain"

func CreateRawCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "raw <get/put/patch/post/delete> v1/url",
		Example: "corectl raw get v1/items --query name=ImportantApp",
		Short:   "Send Http API Request to Qlik Sense Cloud",
		Long:    "Send Http API Request to Qlik Sense Cloud. Query parameters are specified using the --query flag, a body can be specified using one of the body flags (body, body-file or body-values)",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			comm := boot.NewCommunicator(cmd)
			restCaller := comm.RestCaller()

			//Gather parameters
			method := strings.ToUpper(args[0])
			url := args[1]
			queryParams := comm.GetStringMap("query")
			outputFilePath := comm.GetString("output-file")
			bodyFile := comm.GetString("body-file")
			bodyString := comm.GetString("body")
			bodyParams := comm.GetStringMap("body-values")

			//Output the results to either a file or system out
			var out io.Writer
			if outputFilePath != "" {
				file, err := os.Create(outputFilePath)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				out = file
			} else {
				out = cmd.OutOrStdout()
			}

			//Compile the body
			body, mimeType, err := getBodyFromFlags(bodyFile, bodyString, bodyParams)
			if err != nil {
				log.Fatal(err)
			}

			makeRawApiCall(restCaller, method, url, queryParams, body, mimeType, out)
		},
	}
	command.PersistentFlags().StringToStringP("query", "", nil, "Query parameters specified as key=value pairs separated by comma")
	command.PersistentFlags().StringToStringP("body-values", "", nil, "A set of key=value pairs that well be compiled into a json object. A dot (.) inside the key is used to traverse into nested objects. "+
		"The key suffixes :bool or :number can be appended to the key to inject the value into the json structure as boolean or number respectively.")
	command.PersistentFlags().String("body", "", "The content of the body as a string")
	command.PersistentFlags().String("body-file", "", "A file path pointing to a file containing the body of the http request")
	command.PersistentFlags().String("output-file", "", "A file path pointing to where the response body shoule be written")
	return command
}

//makeRawApiCall makes the actual call and writes the output to the supplied io.Writer
func makeRawApiCall(restCaller *rest.RestCaller, method string, url string, queryParams map[string]string, body io.ReadCloser, mimeType string, out io.Writer) (*http.Response, error) {
	defer body.Close()
	//Make the request
	fullUrl := restCaller.CreateUrl(url, queryParams)
	req, err := http.NewRequest(method, fullUrl.String(), body)
	req.Header.Add("Content-Type", mimeType)
	res, err := restCaller.CallRawAndFollowRedirect(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	io.Copy(out, res.Body)
	return res, err
}

// getBodyFromFlags returns a ReadCloser that represents the body regardless of the parameters used
func getBodyFromFlags(bodyFile string, bodyString string, bodyParams map[string]string) (body io.ReadCloser, mimeType string, error error) {
	if multipleBodyParamsSet(bodyString, bodyFile, bodyParams) {
		return nil, "", errors.New("Only one of body-path, body-string, and body-params can be used at one time")
	}
	if bodyString != "" {
		mimeType := detectStringMimeType(bodyString)
		return ioutil.NopCloser(strings.NewReader(bodyString)), mimeType, nil
	}
	if bodyFile != "" {
		mimeType := detectFileMimeType(bodyFile)
		file, err := os.Open(bodyFile)
		return file, mimeType, err
	}
	if bodyParams != nil && len(bodyParams) > 0 {
		body, err := buildBodyFromParams(bodyParams)
		return body, "application/json", err
	}
	return ioutil.NopCloser(strings.NewReader("")), "", nil
}

func buildBodyFromParams(params map[string]string) (io.ReadCloser, error) {
	rootJSonNode := make(map[string]interface{})
	for key_, value := range params {
		keyNameAndType := strings.Split(key_, ":")
		key := keyNameAndType[0]
		valueType := "string"
		if len(keyNameAndType) > 2 {
			return nil, errors.New("invalid key format: " + key_)
		}
		if len(keyNameAndType) == 2 {
			valueType = keyNameAndType[1]
		}
		keyParts := strings.Split(key, ".")
		for _, part := range keyParts {
			if part == "" {
				return nil, errors.New("invalid key format: " + key_)
			}
		}
		jsonNode := rootJSonNode
		for i := 0; i < len(keyParts)-1; i++ {
			keyPart := keyParts[i]
			if jsonNode[keyPart] == nil {
				newNode := make(map[string]interface{})
				jsonNode[keyPart] = newNode
				jsonNode = newNode
			}
		}
		lastKey := keyParts[len(keyParts)-1]
		switch valueType {
		case "int", "integer", "number":
			intValue, intErr := strconv.Atoi(value)
			if intErr != nil {
				floatValue, floatErr := strconv.ParseFloat(value, 64)
				if floatErr != nil {
					if floatErr != nil {
						return nil, floatErr
					}
				}
				jsonNode[lastKey] = floatValue
			}
			jsonNode[lastKey] = intValue

		case "bool", "boolean":
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			jsonNode[lastKey] = boolValue
		case "string":
			jsonNode[lastKey] = value
		default:
			return nil, errors.New("invalid key format: " + key_)
		}

	}
	jsonBytes, err := json.Marshal(rootJSonNode)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(jsonBytes)), nil
}

//detectFileMimeType detects the content-type of the file by checking against various formats
func detectFileMimeType(filePath string) string {
	isJsonFile := strings.HasSuffix(strings.ToLower(filePath), ".json")
	if isJsonFile {
		return applicationJson
	}
	isQvfFile := strings.HasSuffix(strings.ToLower(filePath), ".qvf")
	if isQvfFile {
		return binaryOctetStream
	}

	//Open the file and fetch info
	file, err := os.Open(filePath)
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return binaryOctetStream
	}

	//Determine the buffer size to use
	bufSize := fileInfo.Size()
	if bufSize > maxJsonFileSize {
		bufSize = maxJsonFileSize
	}
	buffer := make([]byte, bufSize)

	// Read the file contents
	_, err = file.Read(buffer)
	// Any errors imply that this is an unknown kind of file
	if err != nil {
		return binaryOctetStream
	}
	// Check the default mimne type
	mimeType := http.DetectContentType(buffer)
	// See if the text/plain detected can be narrowed down to application/json by parsing the content
	if strings.Contains(mimeType, textPlain) && json.Valid(buffer) {
		return applicationJson
	}
	// Translate to another octet-stream content type used by Qlik Sense
	if mimeType == applicationOctetStream {
		return binaryOctetStream
	}
	return mimeType
}

//detectStringMimeType returns either application/json or text/plain
func detectStringMimeType(buffer string) string {
	if json.Valid([]byte(buffer)) {
		return applicationJson
	}
	return textPlain
}

//multipleBodyParamsSet check that only one kind of body parameter is used
func multipleBodyParamsSet(bodyString string, bodyFile string, bodyParams map[string]string) bool {
	bodyCount := 0
	if bodyString != "" {
		bodyCount++
	}
	if bodyFile != "" {
		bodyCount++
	}
	if bodyParams != nil && len(bodyParams) > 0 {
		bodyCount++
	}
	return bodyCount > 1
}
