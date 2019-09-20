package internal

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

// Eval builds a straight table  hypercube based on the supplied argument, evaluates it and prints the result to system out.
func Eval(ctx context.Context, doc *enigma.Doc, args []string) {

	ensureModelExists(ctx, doc)

	measures, dims := argumentsToMeasuresAndDims(args)
	object, _ := doc.CreateSessionObject(ctx, &enigma.GenericObjectProperties{
		Info: &enigma.NxInfo{
			Type: "my-straight-hypercube",
		},
		HyperCubeDef: &enigma.HyperCubeDef{
			Dimensions: createDimensions(dims),
			Measures:   createMeasures(measures),
			InitialDataFetch: []*enigma.NxPage{{
				Height: 20,
				Width:  50,
			}},
		},
	})
	layout, err := object.GetLayout(ctx)

	if err != nil {
		log.Fatalf("could not get hypercube layout: %s\n", err)
	}

	// If the dimension info contains an error element the expression failed to evaluate
	if len(layout.HyperCube.DimensionInfo) != 0 && layout.HyperCube.DimensionInfo[0].Error != nil {
		log.Fatalf("could not evaluate expression, error returned code: %d\n", layout.HyperCube.DimensionInfo[0].Error.ErrorCode)
	}

	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetAutoFormatHeaders(false)

	headers := append(dims, measures...)
	writer.SetHeader(headers)
	writer.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, page := range layout.HyperCube.DataPages {
		for _, row := range page.Matrix {
			data := []string{}
			for _, cell := range row {
				data = append(data, cell.Text)
			}
			writer.Append(data)
		}
	}
	writer.Render()
}

func argumentsToMeasuresAndDims(args []string) ([]string, []string) {
	var (
		measures   = []string{}
		dimensions = []string{}
		foundDims  = false
	)

	for _, arg := range args {
		if arg == "by" {
			foundDims = true
		} else if foundDims {
			dimensions = append(dimensions, arg)
		} else {
			measures = append(measures, arg)
		}
	}

	return measures, dimensions
}
