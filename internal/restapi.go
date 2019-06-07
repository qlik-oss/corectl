package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
)

// ReadRestMetadata fetches metadata from the rest api.
func ReadRestMetadata(url string, headers http.Header) (*RestMetadata, error) {
	if url == "" {
		return nil, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	req.Header = headers
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, nil
	}
	data, _ := ioutil.ReadAll(response.Body)
	result := &RestMetadata{}
	json.Unmarshal(data, result)
	return result, nil
}

func ImportApp(appPath, engine string, headers http.Header) string {
	url, err := neturl.Parse(buildRestBaseURL(engine))
	url.Path = "/v1/apps/import"
	values := neturl.Values{} //name and mode don't seem to work
	url.RawQuery = values.Encode()
	file, err := os.Open(appPath)
	if err != nil {
		FatalError("could not open file: ", appPath)
	}
	defer file.Close()
	req, err := http.NewRequest("POST", url.String(), file)
	if err != nil {
		FatalError(err)
	}
	req.Header = headers
	req.Header.Add("Content-Type", "binary/octet-stream")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		FatalError(err)
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		FatalErrorf("could not create app: got status %d with message %s",
			response.StatusCode, string(data))
	}
	appInfo := &RestNxApp{}
	json.Unmarshal(data, appInfo)
	appID := appInfo.Attributes["id"]
	setAppIDToKnownApps(engine, appID, appID, false)
	return appID
}

func toFieldMetadataMap(fields []*RestFieldMetadata) map[string]*RestFieldMetadata {
	result := make(map[string]*RestFieldMetadata)
	for _, field := range fields {
		result[field.Name] = field
	}
	return result
}

func toTableMetadataMap(tables []*RestTableMetadata) map[string]*RestTableMetadata {
	result := make(map[string]*RestTableMetadata)
	for _, table := range tables {
		result[table.Name] = table
	}
	return result
}

//RestMetadata defines all available info from the metdata rest api.
type RestMetadata struct {
	StaticByteSize int                  `json:"static_byte_size,omitempty"`
	Fields         []*RestFieldMetadata `json:"fields"`
	Tables         []*RestTableMetadata `json:"tables"`
}

func (m *RestMetadata) tableByName(name string) *RestTableMetadata {
	if m != nil {
		for _, table := range m.Tables {
			if table.Name == name {
				return table
			}
		}
	}
	return nil
}

func (m *RestMetadata) fieldByName(name string) *RestFieldMetadata {
	if m != nil {
		for _, field := range m.Fields {
			if field.Name == name {
				return field
			}
		}
	}
	return nil
}

//RestTableMetadata defines all available info about a table.
type RestTableMetadata struct {
	// Name of the table.
	Name string `json:"name"`
	// If set to true, it means that the table is a system table. The default value is false.
	IsSystem bool `json:"is_system"`
	// If set to true, it means that the table is a semantic. The default value is false.
	IsSemantic bool `json:"is_semantic"`
	// If set to true, it means that the table is loose due to circular connection. The default value is false.
	IsLoose bool `json:"is_loose"`
	// of rows.
	NoOfRows int `json:"no_of_rows"`
	// of fields.
	NoOfFields int `json:"no_of_fields"`
	// of key fields.
	NoOfKeyFields int `json:"no_of_key_fields"`
	// Table comment.
	Comment string `json:"commen"`
	// RAM memory used in bytes.
	ByteSize int `json:"byte_size"`
}

//RestFieldMetadata defines all available info about a a field
type RestFieldMetadata struct {
	// Name of the field.
	Name string `json:"name"`
	// No List of table names.
	SrcTables []string `json:"src_tables"`
	// If set to true, it means that the field is a system field. The default value is false.
	IsSystem bool `json:"is_system"`
	// If set to true, it means that the field is hidden. The default value is false.
	IsHidden bool `json:"is_hidden"`
	// If set to true, it means that the field is a semantic. The default value is false.
	IsSemantic bool `json:"is_semantic"`
	// If set to true, only distinct field values are shown. The default value is false.
	DistinctOnly bool `json:"distinct_only"`
	// Number of distinct field values.
	Cardinal int `json:"cardinal"`
	// Total number of field values.
	TotalCount int `json:"total_count"`
	// If set to true, it means that the field is locked. The default value is false.
	IsLocked bool `json:"is_locked"`
	// If set to true, it means that the field has one and only one selection (not 0 and not more than 1). If this property is set to true, the field cannot be cleared anymore and no more selections can be performed in that field. The default value is false.
	AlwaysOneSelected bool `json:"always_one_selected"`
	// Is set to true if the value is a numeric. The default value is false.
	IsNumeric bool `json:"is_numeric"`
	// Field comment.
	Comment string `json:"comment"`
	// No Gives information on a field. For example, it can return the type of the field. Examples: key, text, ASCII.
	Tags []string `json:"tags"`
	// Static RAM memory used in bytes.
	ByteSize int `json:"byte_size"`
}

type RestNxApp struct {
	Attributes map[string]string `json:"attributes"`
}
