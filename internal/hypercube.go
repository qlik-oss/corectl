package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/qlik-oss/enigma-go"
)

func applySelection(ctx context.Context, doc *enigma.Doc, fieldName string, value string) {
	field, err := doc.GetField(ctx, fieldName, "")
	if err != nil {
		fmt.Println("Could not find field", err)
		os.Exit(1)
	}
	field.Clear(ctx)
	if value != "" {
		field.Select(ctx, value, false, 0)
	}
}

func createHypercube(ctx context.Context, doc *enigma.Doc, dimensions []*enigma.NxDimension, measures []*enigma.NxMeasure, sortOrder []int) *enigma.GenericObject {

	object, _ := doc.CreateSessionObject(ctx, &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "my-straight-hypercube",
		},
		HyperCubeDef: &enigma.HyperCubeDef{
			Dimensions: dimensions,
			Measures:   measures,

			InitialDataFetch: []*enigma.NxPage{{
				Height: 20,
				Width:  50,
			}},
			InterColumnSortOrder: sortOrder,
		},
	})

	return object
}

func createMeasures(measures []string) []*enigma.NxMeasure {
	result := make([]*enigma.NxMeasure, len(measures))
	for i, measure := range measures {
		result[i] = createMeasure(measure)
	}
	return result
}

func createMeasure(measure string) *enigma.NxMeasure {
	return &enigma.NxMeasure{
		Def: &enigma.NxInlineMeasureDef{
			Def: measure,
		},
		SortBy: &enigma.SortCriteria{SortByFrequency: 1},
	}
}

func createMeasureSortNumeric(measure string, sort *enigma.SortCriteria) *enigma.NxMeasure {
	return &enigma.NxMeasure{
		Def: &enigma.NxInlineMeasureDef{
			Def: measure,
		},
		SortBy: sort,
	}
}

func createDimensions(dimensions []string) []*enigma.NxDimension {
	result := make([]*enigma.NxDimension, len(dimensions))
	for i, dimension := range dimensions {
		result[i] = createDimension(dimension)
	}
	return result
}

func createDimension(dimension string) *enigma.NxDimension {
	result := &enigma.NxDimension{
		Def: &enigma.NxInlineDimensionDef{
			FieldDefs: []string{dimension},
		},
	}
	return result
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
				Height: 20,
				Width:  50,
			}},
			Mode:         "P",
			NoOfLeftDims: &noOfLeftDims,
		},
	})

	return object
}
