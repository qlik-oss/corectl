package internal

import (
	"context"
	"fmt"

	"github.com/qlik-oss/enigma-go"
)

// PrintFieldValues prints the first few rows of a field to system out.
func PrintFieldValues(ctx context.Context, doc *enigma.Doc, fieldName string) {
	ensureModelExists(ctx, doc)
	fmt.Print(getFieldContentAsTable(ctx, doc, fieldName, 100))
}

func getFieldContentAsString(ctx context.Context, doc *enigma.Doc, fieldName string, length int) string {
	content := getFieldContent(ctx, doc, fieldName, length/2)
	if len(content) > 0 {
		firstItem := content[0]
		if len(content) > 0 && len(firstItem) > length {
			return firstItem[0:length-3] + "..."
		}
		var str string
		for _, item := range content {
			if len(str)+len(item)+2 < length {
				if str != "" {
					str += ", "
				}
				str += item
			} else {
				break
			}
		}
		return str
	}
	return ""
}

func getFieldContentAsTable(ctx context.Context, doc *enigma.Doc, fieldName string, length int) string {
	content := getFieldContent(ctx, doc, fieldName, length)
	if len(content) > 0 {
		firstItem := content[0]
		if len(content) > 0 && len(firstItem) > length {
			return firstItem[0:length-3] + "..."
		}
		var str string
		for _, item := range content {
			str += item + "\n"
		}
		return str
	}
	return ""
}

func getFieldContent(ctx context.Context, doc *enigma.Doc, fieldName string, count int) []string {

	object, _ := doc.CreateSessionObject(ctx, &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "my-straight-hypercube",
		},
		ListObjectDef: &enigma.ListObjectDef{
			Def: &enigma.NxInlineDimensionDef{
				FieldDefs:     []string{fieldName},
				SortCriterias: []*enigma.SortCriteria{{SortByFrequency: 1}},
			},
			ShowAlternatives: true,
			FrequencyMode:    "NX_FREQUENCY_VALUE",
			InitialDataFetch: []*enigma.NxPage{
				{
					Top:    0,
					Height: count,
					Left:   0,
					Width:  1,
				},
			},
		},
	})

	layout, err := object.GetLayout(ctx)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var result []string
	// Get hypercube layout

	for _, page := range layout.ListObject.DataPages {
		for _, row := range page.Matrix {
			for _, cell := range row {
				if cell.Frequency != "1" && cell.Frequency != "" {
					result = append(result, cell.Text+"("+cell.Frequency+"x)")
				} else {
					result = append(result, cell.Text)
				}

			}
		}
	}
	return result
}
