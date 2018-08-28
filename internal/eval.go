package internal

import (
	"context"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/qlik-oss/enigma-go"
	"strings"
)

func Eval(ctx context.Context, doc *enigma.Doc, args []string) {
	measures, dims := argumentsToMeasuresAndDims(args)
	object, _ := doc.CreateObject(ctx, &enigma.GenericObjectProperties{
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
	fmt.Println("---------- " + strings.Join(args, " ") + " ----------")
	grid := tm.NewTable(0, 10, 3, ' ', 0)
	layout, _ := object.GetLayout(ctx)
	fmt.Print(grid, strings.Join(dims, "\t"))
	fmt.Print(grid, "\t")
	fmt.Println(grid, strings.Join(measures, "\t"))
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
	tm.Println(grid)
	tm.Flush()
}

func argumentsToMeasuresAndDims(args []string) ([]string, []string) {
	var (
		dims      []string
		measures  []string
		tempArray []string
	)
	for _, arg := range args {
		if arg != "by" {
			tempArray = append(tempArray, arg)
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
