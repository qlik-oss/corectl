package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/qlik-oss/enigma-go"
)

type ModelMetadata struct {
	Tables                     []*enigma.TableRecord
	SourceKeys                 []*enigma.SourceKeyRecord
	Metadata                   *Metadata
	TableMetadataByName        map[string]*TableMetadata
	FieldMetadataByName        map[string]*FieldMetadata
	SystemTableLayout          *enigma.GenericObjectLayout
	FieldsInTable              map[string]string
	FieldNames                 []string
	TableNames                 []string
	FieldSourceTableInfoByName map[string][]FieldSourceTableInfo
	FieldInTableDataByNames    map[string]map[string]*enigma.FieldInTableData
	SampleContentByFieldName   map[string]string
}

func GetModelMetadata(ctx context.Context, doc *enigma.Doc, metaURL string, keyOnly bool) *ModelMetadata {
	tables, sourceKeys, err := doc.GetTablesAndKeys(ctx, &enigma.Size{}, &enigma.Size{}, 0, false, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	metadata, err := ReadMetadata(metaURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(tables) > 0 && len(metadata.Tables) == 0 {
		fmt.Println("Could not load metadata")
		os.Exit(1)
	}
	systemTableObject := createSystemTableHypercube(ctx, doc)
	systemTableLayout, _ := systemTableObject.GetLayout(ctx)
	fieldsInTable := tableRecordsToMap(tables)

	fieldInTableDataByNames := tableRecordsToMapMap(tables)

	tableNames, fieldNames, tableFieldInfoByFieldName := systemTableToSystemMap(systemTableLayout, fieldInTableDataByNames)
	return &ModelMetadata{
		Tables:                     tables,
		SourceKeys:                 sourceKeys,
		TableMetadataByName:        ToTableMetadataMap(metadata.Tables),
		FieldMetadataByName:        ToFieldMetadataMap(metadata.Fields),
		Metadata:                   metadata,
		SystemTableLayout:          systemTableLayout,
		FieldsInTable:              fieldsInTable,
		FieldNames:                 fieldNames,
		TableNames:                 tableNames,
		FieldSourceTableInfoByName: tableFieldInfoByFieldName,
		FieldInTableDataByNames:    fieldInTableDataByNames,
		SampleContentByFieldName:   buildSampleContent(ctx, doc, fieldNames),
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

type FieldSourceTableInfo struct {
	RowCount string
	KeyType  string
}
type SystemTableInfo struct {
	TableNames                []string
	TableFieldInfoByFieldName map[string][]FieldSourceTableInfo
}

func systemTableToSystemMap(systemTableLayout *enigma.GenericObjectLayout, fieldInfo map[string]map[string]*enigma.FieldInTableData) ([]string, []string, map[string][]FieldSourceTableInfo) {

	page := systemTableLayout.HyperCube.PivotDataPages[0]

	tableNames := make([]string, len(page.Top))
	fieldNames := make([]string, len(page.Left))
	tableFieldInfoByFieldName := make(map[string][]FieldSourceTableInfo)

	for x, tableName := range page.Top {
		tableNames[x] = tableName.Text
	}

	for y, row := range page.Data {
		fieldName := page.Left[y].Text
		fieldNames[y] = fieldName
		fieldSourceTableInfos := make([]FieldSourceTableInfo, len(page.Top))

		for x := range row {

			tableName := tableNames[x]
			fieldDetails := fieldInfo[tableName][fieldName]
			if fieldDetails != nil {

				info := fmt.Sprintf("%d/%d", fieldDetails.NTotalDistinctValues, fieldDetails.NNonNulls)
				if fieldDetails.NRows > fieldDetails.NNonNulls {
					info += fmt.Sprintf("+%d", fieldDetails.NRows-fieldDetails.NNonNulls)
				}
				if fieldDetails.KeyType == "NOT_KEY" {
					info += ""
				} else if fieldDetails.KeyType == "ANY_KEY" {
					info += "*"
				} else if fieldDetails.KeyType == "PRIMARY_KEY" {
					info += "**"
				} else if fieldDetails.KeyType == "PERFECT_KEY" {
					info += "***"
				} else {
					info += "?"
				}
				fieldSourceTableInfos[x].RowCount = info
			} else {
				fieldSourceTableInfos[x].RowCount = ""
			}

		}
		tableFieldInfoByFieldName[fieldName] = fieldSourceTableInfos
	}

	return tableNames, fieldNames, tableFieldInfoByFieldName
}
