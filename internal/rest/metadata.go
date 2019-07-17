package rest

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
)

// ReadRestMetadata fetches the metadata for the specified app
func ReadRestMetadata(appID string, engine *neturl.URL, headers http.Header, certs *tls.Config) (*RestMetadata, error) {
	url := CreateBaseURL(*engine)
	url.Path = fmt.Sprintf("/v1/apps/%s/data/metadata", adaptAppID(appID))
	req := &http.Request{
		Method: "GET",
		URL:    url,
		Header: headers,
	}
	result := &RestMetadata{}
	statusCodes := &map[int]bool{
		200: true,
	}
	err := Call(req, certs, result, statusCodes, json.Unmarshal)
	if err != nil {
		return nil, fmt.Errorf("could not get rest metadata: %s", err.Error())
	}
	return result, nil
}

type RestMetadata struct {
	StaticByteSize int                  `json:"static_byte_size,omitempty"`
	Fields         []*RestFieldMetadata `json:"fields"`
	Tables         []*RestTableMetadata `json:"tables"`
}

func (m *RestMetadata) TableByName(name string) *RestTableMetadata {
	if m != nil {
		for _, table := range m.Tables {
			if table.Name == name {
				return table
			}
		}
	}
	return nil
}

func (m *RestMetadata) FieldByName(name string) *RestFieldMetadata {
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
