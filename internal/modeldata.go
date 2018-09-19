package internal

import (
	"context"
	"fmt"
	"github.com/qlik-oss/enigma-go"
	"os"
)

// ModelMetadata defines all available metadata around the data model.
type ModelMetadata struct {
	Tables                   []*TableModel
	Fields                   []*FieldModel
	SourceKeys               []*enigma.SourceKeyRecord
	RestMetadata             *RestMetadata
	RestTableMetadataByName  map[string]*RestTableMetadata
	RestFieldMetadataByName  map[string]*RestFieldMetadata
	FieldsInTableTexts       map[string]string
	SampleContentByFieldName map[string]string
}

// TableModel represents one table in the data model. It contains information from both the QIX and Rest apis
type TableModel struct {
	*enigma.TableRecord
	RestMetadata *RestTableMetadata
}

// FieldModel represents one field in the data model. It contains information from both the QIX and Rest apis.
// It also contains an array (with compatible ordering with the main Table model ordering) with field per source table info.
type FieldModel struct {
	*enigma.FieldDescription
	RestMetadata *RestFieldMetadata
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

func createFieldModels(ctx context.Context, doc *enigma.Doc, fieldNames []string, restMetadata *RestMetadata) []*FieldModel {
	result := make([]*FieldModel, len(fieldNames))
	resultChannel := make(chan *enigma.FieldDescription, len(fieldNames))
	for _, fieldName := range fieldNames {
		fieldDescr, err := doc.GetFieldDescription(ctx, fieldName)
		if err != nil {
			fmt.Println("Unexpected error", err)
			os.Exit(1)
		}
		resultChannel <- fieldDescr
	}
	for i, fieldName := range fieldNames {
		fieldDescr := <-resultChannel

		result[i] = &FieldModel{
			FieldDescription: fieldDescr,
			RestMetadata:     restMetadata.fieldByName(fieldName),
		}
	}

	return result
}

func createTableModels(ctx context.Context, doc *enigma.Doc, tableNames []string, tableRecords []*enigma.TableRecord, restMetadata *RestMetadata) []*TableModel {
	tableRecordMap := make(map[string]*enigma.TableRecord)
	for _, tableRecord := range tableRecords {
		tableRecordMap[tableRecord.Name] = tableRecord
	}

	tableModels := make([]*TableModel, len(tableNames))
	for i, tableName := range tableNames {
		tableModels[i] = &TableModel{TableRecord: tableRecordMap[tableName], RestMetadata: restMetadata.tableByName(tableName)}
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
func GetModelMetadata(ctx context.Context, doc *enigma.Doc, metaURL string, keyOnly bool) *ModelMetadata {
	tables, sourceKeys, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(tables) == 0 {
		fmt.Println("The data model is empty.")
		os.Exit(1)
	}
	restMetadata, err := ReadRestMetadata(metaURL)

	if len(tables) > 0 && restMetadata == nil {
		fmt.Println("No REST metadata available.")
	}
	tableNames, fieldNames := getSortedTableNamesAndFieldsNames(ctx, doc, err, tables)

	fieldModels := createFieldModels(ctx, doc, fieldNames, restMetadata)
	tableModels := createTableModels(ctx, doc, tableNames, tables, restMetadata)
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

func buildSampleContent(ctx context.Context, doc *enigma.Doc, fieldNames []string) map[string]string {
	var sampleContentByFieldName map[string]string
	sampleContentByFieldName = make(map[string]string)
	for _, fieldName := range fieldNames {
		sampleContentByFieldName[fieldName] = getFieldContentAsString(ctx, doc, fieldName, 40)
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

// Exists if there  is no data model
func ensureModelExists(ctx context.Context, doc *enigma.Doc) {
	tables, _, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(tables) == 0 {
		fmt.Println("The data model is empty.")
		os.Exit(1)
	}
}

func DataModelTableCount(ctx context.Context, doc *enigma.Doc) int {
	tables, _, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		return 0
	}
	return len(tables)
}
