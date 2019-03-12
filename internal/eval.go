package internal

import (
	"context"
	"fmt"
	"strings"

	tm "github.com/buger/goterm"
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
	grid := tm.NewTable(0, 10, 3, ' ', 0)
	layout, err := object.GetLayout(ctx)

	if err != nil {
		Logger.Fatal("Failed to get hypercube layout: ", err)
	}

	// If the dimension info contains an error element the expression failed to evaluate
	if len(layout.HyperCube.DimensionInfo) != 0 && layout.HyperCube.DimensionInfo[0].Error != nil {
		Logger.Fatalf("Failed to evaluate expression with error code: %d", layout.HyperCube.DimensionInfo[0].Error.ErrorCode)
	}

	fmt.Fprintf(grid, strings.Join(dims, "\t"))
	fmt.Fprintf(grid, "\t")
	fmt.Fprintf(grid, strings.Join(measures, "\t"))
	fmt.Fprintf(grid, "\n")
	// Get hypercube layout
	for _, page := range layout.HyperCube.DataPages {
		for _, row := range page.Matrix {
			for r, cell := range row {
				if r < len(row)-1 {
					fmt.Fprintf(grid, "%s\t", cell.Text)
				} else {
					fmt.Fprintf(grid, "%s\n", cell.Text)
				}
			}
		}
	}
	fmt.Print(grid)
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
