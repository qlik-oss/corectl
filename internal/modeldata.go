package internal

import (
	"context"
	"crypto/tls"
	"net/http"
	neturl "net/url"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/corectl/internal/rest"
	"github.com/qlik-oss/enigma-go"
)

// ModelMetadata defines all available metadata around the data model.
type ModelMetadata struct {
	Tables                   []*TableModel
	Fields                   []*FieldModel
	SourceKeys               []*enigma.SourceKeyRecord
	RestMetadata             *rest.RestMetadata
	RestTableMetadataByName  map[string]*rest.RestTableMetadata
	RestFieldMetadataByName  map[string]*rest.RestFieldMetadata
	FieldsInTableTexts       map[string]string
	SampleContentByFieldName map[string]string
}

// TableModel represents one table in the data model. It contains information from both the QIX and Rest apis
type TableModel struct {
	*enigma.TableRecord
	RestMetadata *rest.RestTableMetadata
}

// FieldModel represents one field in the data model. It contains information from both the QIX and Rest apis.
// It also contains an array (with compatible ordering with the main Table model ordering) with field per source table info.
type FieldModel struct {
	*enigma.FieldDescription
	RestMetadata *rest.RestFieldMetadata
	//Sparse array with information about the source tables.
	FieldInTable []*enigma.FieldInTableData
}

// MemUsage returns the memory usage by the table if that information is available, "N/A" otherwise.
func (t *TableModel) MemUsage() string {
	if t.RestMetadata != nil {
		return FormatBytes(t.RestMetadata.ByteSize)
	}
	return "N/A"
}

// MemUsage returns the static memory usage of the whole data model if that information is available, "N/A" otherwise.
func (m *ModelMetadata) MemUsage() string {
	if m.RestMetadata != nil {
		return FormatBytes(m.RestMetadata.StaticByteSize)
	}
	return "N/A"
}

// MemUsage returns the memory usage of the field.
func (f *FieldModel) MemUsage() string {
	return FormatBytes(f.ByteSize)
}

func (m *ModelMetadata) tableByName(name string) *TableModel {
	if m != nil {
		for _, table := range m.Tables {
			if table.Name == name {
				return table
			}
		}
	}
	return nil
}

func createFieldModels(ctx context.Context, doc *enigma.Doc, fieldNames []string, restMetadata *rest.RestMetadata) []*FieldModel {
	result := make([]*FieldModel, len(fieldNames))

	type GetFieldDescriptionResultEntry struct {
		index  int
		result *enigma.FieldDescription
	}
	waitChannel := make(chan GetFieldDescriptionResultEntry)
	defer close(waitChannel)

	for i, fieldName := range fieldNames {
		result[i] = &FieldModel{
			FieldDescription: nil, //Fill in later
			RestMetadata:     restMetadata.FieldByName(fieldName),
		}
		//Run field description fetching in parallel threads
		go func(index int, fieldName string) {
			fieldDescr, err := doc.GetFieldDescription(ctx, fieldName)
			if err != nil {
				log.Fatalf("could not retrieve field description for '%s': %s\n", fieldName, err)
			}
			item := GetFieldDescriptionResultEntry{index: index, result: fieldDescr}
			waitChannel <- item
		}(i, fieldName)
	}
	for range fieldNames {
		item := <-waitChannel
		result[item.index].FieldDescription = item.result
	}

	return result
}

func createTableModels(ctx context.Context, doc *enigma.Doc, tableRecords []*enigma.TableRecord, restMetadata *rest.RestMetadata) []*TableModel {
	tableModels := make([]*TableModel, len(tableRecords))
	for i, tableRecord := range tableRecords {
		tableModels[i] = &TableModel{TableRecord: tableRecord, RestMetadata: restMetadata.TableByName(tableRecord.Name)}
	}
	return tableModels
}

func addTableFieldCellCrossReferences(fields []*FieldModel, tables []*TableModel) {
	for _, field := range fields {
		field.FieldInTable = make([]*enigma.FieldInTableData, len(tables))
		for tableIndex, table := range tables {
			for _, fieldInTableData := range table.Fields {
				if fieldInTableData.Name == field.Name {
					field.FieldInTable[tableIndex] = fieldInTableData
				}
			}

		}
	}
}

// GetModelMetadata retrives all available metadata about the app
func GetModelMetadata(ctx context.Context, doc *enigma.Doc, appID string, engine *neturl.URL, headers http.Header, tlsClientConfig *tls.Config, keyOnly bool) *ModelMetadata {
	tables, sourceKeys, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		log.Fatalf("could not retrieve tables and keys: %s\n", err)
	}
	if len(tables) == 0 {
		log.Fatalf("the data model is empty\n")
	}
	restMetadata, err := rest.ReadRestMetadata(appID, engine, headers, tlsClientConfig)

	if len(tables) > 0 && restMetadata == nil {
		log.Infoln("No REST metadata available.")
	}
	fieldNames := getSortedFieldsNames(ctx, doc, err)

	fieldModels := createFieldModels(ctx, doc, fieldNames, restMetadata)
	if keyOnly {
		fieldModels = filterKeyFields(fieldModels)
	}
	tableModels := createTableModels(ctx, doc, tables, restMetadata)
	fieldsInTableTexts := tableRecordsToMap(tables)
	addTableFieldCellCrossReferences(fieldModels, tableModels)
	return &ModelMetadata{
		Tables:                   tableModels,
		Fields:                   fieldModels,
		SourceKeys:               sourceKeys,
		RestMetadata:             restMetadata,
		FieldsInTableTexts:       fieldsInTableTexts,
		SampleContentByFieldName: buildSampleContent(ctx, doc, fieldNames),
	}
}

func filterKeyFields(fields []*FieldModel) []*FieldModel {
	filteredFields := []*FieldModel{}
	for _, field := range fields {
		if isKey(field) {
			filteredFields = append(filteredFields, field)
		}
	}
	return filteredFields
}

func isKey(field *FieldModel) bool {
	for _, tag := range field.Tags {
		if tag == "$key" {
			return true
		}
	}
	return false
}

func buildSampleContent(ctx context.Context, doc *enigma.Doc, fieldNames []string) map[string]string {
	var sampleContentByFieldName map[string]string
	sampleContentByFieldName = make(map[string]string)

	type GetFieldContentAsStringResultEntry struct {
		fieldName     string
		sampleContent string
	}
	waitChannel := make(chan GetFieldContentAsStringResultEntry)
	defer close(waitChannel)

	for _, fieldName := range fieldNames {
		go func(fieldName string) {
			sampleContent := getFieldContentAsString(ctx, doc, fieldName, 40)
			waitChannel <- GetFieldContentAsStringResultEntry{fieldName: fieldName, sampleContent: sampleContent}
		}(fieldName)
	}

	for range fieldNames {
		item := <-waitChannel
		sampleContentByFieldName[item.fieldName] = item.sampleContent
	}

	return sampleContentByFieldName
}

func tableRecordsToMap(tables []*enigma.TableRecord) map[string]string {
	fieldsInTable := make(map[string]string)
	for _, table := range tables {
		fieldInfo := ""
		for f, field := range table.Fields {
			if f > 0 {
				fieldInfo += ", "
			}
			fieldInfo += field.Name
		}
		fieldsInTable[table.Name] = fieldInfo

	}
	return fieldsInTable
}

func tableRecordsToMapMap(tables []*enigma.TableRecord) map[string]map[string]*enigma.FieldInTableData {
	fieldsInTable := make(map[string]map[string]*enigma.FieldInTableData)
	for _, table := range tables {
		fieldsInTable[table.Name] = make(map[string]*enigma.FieldInTableData)
		for _, field := range table.Fields {
			fieldsInTable[table.Name][field.Name] = field
		}
	}
	return fieldsInTable
}

// FieldSourceTableInfo defines row count and key type for a field
type FieldSourceTableInfo struct {
	CellContent string
	KeyType     string
}

// Exits if there is no data model
func ensureModelExists(ctx context.Context, doc *enigma.Doc) {
	tables, _, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		log.Fatalf("could not retrieve tables and keys: %s\n", err)
	}
	if len(tables) == 0 {
		log.Fatalf("the data model is empty\n")
	}
}
