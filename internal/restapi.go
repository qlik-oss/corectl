package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ReadMetadata fetches metadata from the rest api.
func ReadMetadata(url string) (*Metadata, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Too old app format")
	}
	data, _ := ioutil.ReadAll(response.Body)
	result := &Metadata{}
	json.Unmarshal(data, result)
	return result, nil
}

func toFieldMetadataMap(fields []*FieldMetadata) map[string]*FieldMetadata {
	result := make(map[string]*FieldMetadata)
	for _, field := range fields {
		result[field.Name] = field
	}
	return result
}

func toTableMetadataMap(tables []*TableMetadata) map[string]*TableMetadata {
	result := make(map[string]*TableMetadata)
	for _, table := range tables {
		result[table.Name] = table
	}
	return result
}

//Metadata defines all available info from the metdata rest api.
type Metadata struct {
	StaticByteSize int              `json:"static_byte_size,omitempty"`
	Fields         []*FieldMetadata `json:"fields"`
	Tables         []*TableMetadata `json:"tables"`
}

//TableMetadata defines all available info about a table.
type TableMetadata struct {
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

//FieldMetadata defines all available info about a a field
type FieldMetadata struct {
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
