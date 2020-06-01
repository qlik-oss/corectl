package standard

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/commands/engine"
	"github.com/qlik-oss/corectl/pkg/log"
	"github.com/qlik-oss/corectl/pkg/rest"
	"github.com/spf13/cobra"
)

const maxJsonFileSize = 1024 * 1024
const applicationJson = "application/json"
const binaryOctetStream = "binary/octet-stream"
const applicationOctetStream = "application/octet-stream"
const textPlain = "text/plain"

func CreateRawCommand() *cobra.Command {
	command := engine.WithLocalFlags(&cobra.Command{
		Use:               "raw <get/put/patch/post/delete> v1/url",
		Example:           "corectl raw get v1/items --query name=ImportantApp",
		Short:             "Send Http API Request to Qlik Sense Cloud editions",
		Long:              "Send Http API Request to Qlik Sense Cloud editions. Query parameters are specified using the --query flag, a body can be specified using one of the body flags (body, body-file or body-values)",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: rawCompletion,
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

			//Compile the body
			body, mimeType, err := getBodyFromFlags(bodyFile, bodyString, bodyParams)
			if err != nil {
				log.Fatal(err)
			}

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

			// Setting the "flag" to add raw to the User-Agent
			comm.DynSettings.SetUserAgentComment("raw")

			var filter rest.Filter
			if comm.GetBool("quiet") {
				filter = rest.QuietFilter
			}

			err = restCaller.CallStreaming(method, url, queryParams, mimeType, body, out, filter)

			if err != nil {
				// Cleanup if we're trying to write to a file.
				if file, ok := out.(*os.File); ok {
					// Close since os.Exit doesn't respect deferred functions.
					_ = file.Close()
					// Remove the empty file after a failed call.
					_ = os.Remove(outputFilePath)
				}
				fmt.Fprintln(cmd.ErrOrStderr(), err)
				os.Exit(1)
			}
		},
	}, "quiet")
	command.PersistentFlags().StringToStringP("query", "", nil, "Query parameters specified as key=value pairs separated by comma")
	command.PersistentFlags().StringToStringP("body-values", "", nil, "A set of key=value pairs that well be compiled into a json object. A dot (.) inside the key is used to traverse into nested objects. "+
		"The key suffixes :bool or :number can be appended to the key to inject the value into the json structure as boolean or number respectively.")
	command.PersistentFlags().String("body", "", "The content of the body as a string")
	command.PersistentFlags().String("body-file", "", "A file path pointing to a file containing the body of the http request")
	command.PersistentFlags().String("output-file", "", "A file path pointing to where the response body shoule be written")
	return command
}

// rawCompletion is the completion function for the raw command.
// As the raw command is sort of a "free typing" command, we can only help the user a bit
// on the way.
func rawCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) < 1 {
		return []string{"get", "put", "patch", "post", "delete"}, cobra.ShellCompDirectiveNoFileComp
	}
	if len(args) < 2 {
		return []string{"v1/"}, cobra.ShellCompDirectiveNoSpace
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
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
	isJsonFile := hasFileType(filePath, ".json")
	if isJsonFile {
		return applicationJson
	}
	isQvfFile := hasFileType(filePath, ".qvf")
	isQvwFile := hasFileType(filePath, ".qvw")
	isQvdFile := hasFileType(filePath, ".qvd")
	isQvxFile := hasFileType(filePath, ".qvx")
	if isQvfFile || isQvwFile || isQvdFile || isQvxFile {
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

func hasFileType(filePath string, suffix string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), suffix)
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
