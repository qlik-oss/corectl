package internal

import (
	"context"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

func getSortedFieldsNames(ctx context.Context, doc *enigma.Doc, err error) []string {
	systemTableObject := createSystemTableHypercube(ctx, doc)
	systemTableLayout, err := systemTableObject.GetLayout(ctx)
	if err != nil {
		log.Fatalln("could not fetch system table: ", err)
	}
	fieldNames := layoutToFieldLists(systemTableLayout)
	return fieldNames
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

func layoutToFieldLists(systemTableLayout *enigma.GenericObjectLayout) []string {

	page := systemTableLayout.HyperCube.PivotDataPages[0]
	fieldNames := make([]string, len(page.Left))

	for y := range page.Data {
		fieldName := page.Left[y].Text
		fieldNames[y] = fieldName
	}

	return fieldNames
}
