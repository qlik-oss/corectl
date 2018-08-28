package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func ToFieldMetadataMap(fields []*FieldMetadata) map[string]*FieldMetadata {
	result := make(map[string]*FieldMetadata)
	for _, field := range fields {
		result[field.Name] = field
	}
	return result
}

func ToTableMetadataMap(tables []*TableMetadata) map[string]*TableMetadata {
	result := make(map[string]*TableMetadata)
	for _, table := range tables {
		result[table.Name] = table
	}
	return result
}

type Metadata struct {
	StaticByteSize int              `json:"static_byte_size,omitempty"`
	Fields         []*FieldMetadata `json:"fields"`
	Tables         []*TableMetadata `json:"tables"`
}

type TableMetadata struct {
	// Name of the table.
	Name string `json:"name"`
	// If set to true, it means that the table is a system table. The default value is false.
	Is_system bool `json:"is_system"`
	// If set to true, it means that the table is a semantic. The default value is false.
	Is_semantic bool `json:"is_semantic"`
	// If set to true, it means that the table is loose due to circular connection. The default value is false.
	Is_loose bool `json:"is_loose"`
	// of rows.
	No_of_rows int `json:"no_of_rows"`
	// of fields.
	No_of_fields int `json:"no_of_fields"`
	// of key fields.
	No_of_key_fields int `json:"no_of_key_fields"`
	// Table comment.
	Comment string `json:"commen"`
	// RAM memory used in bytes.
	Byte_size int `json:"byte_size"`
}

type FieldMetadata struct {
	// Name of the field.
	Name string `json:"name"`
	// No List of table names.
	Src_tables []string `json:"src_tables"`
	// If set to true, it means that the field is a system field. The default value is false.
	Is_system bool `json:"is_system"`
	// If set to true, it means that the field is hidden. The default value is false.
	Is_hidden bool `json:"is_hidden"`
	// If set to true, it means that the field is a semantic. The default value is false.
	Is_semantic bool `json:"is_semantic"`
	// If set to true, only distinct field values are shown. The default value is false.
	Distinct_only bool `json:"distinct_only"`
	// Number of distinct field values.
	Cardinal int `json:"cardinal"`
	// Total number of field values.
	Total_count int `json:"total_count"`
	// If set to true, it means that the field is locked. The default value is false.
	Is_locked bool `json:"is_locked"`
	// If set to true, it means that the field has one and only one selection (not 0 and not more than 1). If this property is set to true, the field cannot be cleared anymore and no more selections can be performed in that field. The default value is false.
	Always_one_selected bool `json:"always_one_selected"`
	// Is set to true if the value is a numeric. The default value is false.
	Is_numeric bool `json:"is_numeric"`
	// Field comment.
	Comment string `json:"comment"`
	// No Gives information on a field. For example, it can return the type of the field. Examples: key, text, ASCII.
	Tags []string `json:"tags"`
	// Static RAM memory used in bytes.
	Byte_size int `json:"byte_size"`
}
