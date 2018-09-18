package internal

import (
	"context"
	"fmt"
	"github.com/qlik-oss/enigma-go"
	"os"
)

func getSortedTableNamesAndFieldsNames(ctx context.Context, doc *enigma.Doc, err error, tables []*enigma.TableRecord) ([]string, []string) {
	systemTableObject := createSystemTableHypercube(ctx, doc)
	systemTableLayout, err := systemTableObject.GetLayout(ctx)
	if err != nil {
		fmt.Println("Error when fetching system table:", err)
		os.Exit(1)
	}
	tableNames, fieldNames := layoutToFieldListsAndTableLists(systemTableLayout)
	return tableNames, fieldNames
}

func createSystemTableHypercube(ctx context.Context, doc *enigma.Doc) *enigma.GenericObject {

	noOfLeftDims := 1

	object, _ := doc.CreateSessionObject(ctx, &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "my-pivot-hypercube",
		},
		HyperCubeDef: &enigma.HyperCubeDef{
			Dimensions: []*enigma.NxDimension{
				{
					Def: &enigma.NxInlineDimensionDef{
						FieldDefs:     []string{"=$Field"},
						SortCriterias: []*enigma.SortCriteria{{SortByExpression: -1, Expression: &enigma.ValueExpr{V: "=count($Table)"}}},
					},
				},
				{
					Def: &enigma.NxInlineDimensionDef{
						FieldDefs:     []string{"=$Table"},
						SortCriterias: []*enigma.SortCriteria{{SortByExpression: -1, Expression: &enigma.ValueExpr{V: "=count($Field)"}}},
					},
				},
			},
			Measures: []*enigma.NxMeasure{
				createMeasureSortNumeric("=sum($Rows)", &enigma.SortCriteria{SortByNumeric: -1}),
			},

			InitialDataFetch: []*enigma.NxPage{{
				Height: 1000,
				Width:  1000,
			}},
			Mode:         "P",
			NoOfLeftDims: &noOfLeftDims,
		},
	})

	return object
}

func layoutToFieldListsAndTableLists(systemTableLayout *enigma.GenericObjectLayout) ([]string, []string) {

	page := systemTableLayout.HyperCube.PivotDataPages[0]
	tableNames := make([]string, len(page.Top))
	fieldNames := make([]string, len(page.Left))

	for x, tableName := range page.Top {
		tableNames[x] = tableName.Text
	}

	for y := range page.Data {
		fieldName := page.Left[y].Text
		fieldNames[y] = fieldName
	}

	return tableNames, fieldNames
}
