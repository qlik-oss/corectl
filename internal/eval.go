package internal

import (
	"context"
	"fmt"
	"os"
	"strings"

	tm "github.com/buger/goterm"
	"github.com/qlik-oss/enigma-go"
)

// Eval builds a straight table  hypercube based on the supplied argument, evaluates it and prints the result to system out.
func Eval(ctx context.Context, doc *enigma.Doc, args []string) {
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
		fmt.Println("Failed to get hypercube layout: ", err)
		os.Exit(1)
	}

	// If the dimension info contains an error element the expression failed to evaluate
	if len(layout.HyperCube.DimensionInfo) != 0 && layout.HyperCube.DimensionInfo[0].Error != nil {
		fmt.Println("Failed to evaluate expression with error code:", layout.HyperCube.DimensionInfo[0].Error.ErrorCode)
		os.Exit(1)
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
		dims      []string
		measures  []string
		tempArray []string
	)
	for _, arg := range args {
		if arg != "by" {
			// Skip appending dimension if iterating over all dimensions
			if arg != "*" {
				tempArray = append(tempArray, arg)
			}
		} else {
			//The first set of arguments are treated as measures when we find the "by" keyword
			//Switch to adding dimensions
			measures = tempArray
			tempArray = []string{}
		}
	}
	dims = tempArray
	return measures, dims
}
