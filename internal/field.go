package internal

import (
	"context"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// PrintFieldValues prints the first few rows of a field to system out.
func PrintFieldValues(ctx context.Context, doc *enigma.Doc, fieldName string) {
	ensureModelExists(ctx, doc)
	log.Quiet(getFieldContentAsTable(ctx, doc, fieldName, 100))
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
		log.Errorln(err)
		return []string{}
	}

	var result []string

	// If there are no datapages, it is (probably?) not a field.
	if len(layout.ListObject.DataPages) == 0 {
		log.Fatalf("no field by name '%s'\n", fieldName)
	}

	// Get hypercube layout
	for _, page := range layout.ListObject.DataPages {
		for _, row := range page.Matrix {
			for _, cell := range row {
				text := cell.Text
				// If there is no text, we will represent it as <empty>
				// to align with how catwalk handles such values.
				if len(text) == 0 {
					text = "<empty>"
				}
				if cell.Frequency != "1" && cell.Frequency != "" {
					result = append(result, text+"("+cell.Frequency+"x)")
				} else {
					result = append(result, text)
				}

			}
		}
	}
	return result
}
