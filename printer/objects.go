package printer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/enigma-go"
	"strings"
)

// PrintObjects prints a list of the id and type of all objects in the app
func PrintObjects(allInfos []*enigma.NxInfo) {

	objectTable := tm.NewTable(0, 10, 3, ' ', 0)
	fmt.Fprintf(objectTable, "Id\tType\n")
	for _, info := range allInfos {

		fmt.Fprintf(objectTable, "%s\t%s\n", info.Id, info.Type)
	}
	fmt.Print(objectTable)
}

// PrintObject prints the properties of the object defined by objectID
func PrintObject(state *internal.State, objectID string) {
	object, err := state.Doc.GetObject(state.Ctx, objectID)
	if err != nil {
		internal.FatalError(err)
	}
	properties, err := object.GetPropertiesRaw(state.Ctx)
	if err != nil {
		internal.FatalError(err)
	}
	fmt.Println(prettyJSON(properties))
}

// PrintObjectLayout prints the layout of the object defined by objectID
func PrintObjectLayout(state *internal.State, objectID string) {
	object, err := state.Doc.GetObject(state.Ctx, objectID)
	if err != nil {
		internal.FatalError(err)
	}
	properties, err := object.GetLayoutRaw(state.Ctx)
	if err != nil {
		internal.FatalError(err)
	}
	fmt.Println(prettyJSON(properties))
}
func prettyJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, data, "", "   ")
	return prettyJSON.String()
}

// EvalObject evalutes the data of the object identified by objectID
func EvalObject(ctx context.Context, doc *enigma.Doc, objectID string) {
	object, err := doc.GetObject(ctx, objectID)
	if err != nil {
		internal.FatalError(err)
	}

	layout, err := object.GetLayoutRaw(ctx)
	if err != nil {
		internal.FatalError(err)
	}

	layoutMap := make(map[string]interface{})
	err = json.Unmarshal(layout, &layoutMap)
	if err != nil {
		fmt.Println(err)
	}
	resultCubeMap := make(map[string]*enigma.HyperCube)
	getAllHyperCubes("", layoutMap, resultCubeMap)
	for _, hypercube := range resultCubeMap {
		printHypercube(hypercube)
	}
}

func printHypercube(hypercube *enigma.HyperCube) {
	grid := tm.NewTable(0, 10, 3, ' ', 0)

	for _, dim := range hypercube.DimensionInfo {
		fmt.Fprintf(grid, "%s\t", dim.FallbackTitle)
	}
	for _, mes := range hypercube.MeasureInfo {
		fmt.Fprintf(grid, "%s\t", mes.FallbackTitle)
	}
	fmt.Fprint(grid, "\n")

	for _, page := range hypercube.DataPages {
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
	fmt.Println(grid.String())
}

func getAllHyperCubes(path string, jsonMap map[string]interface{}, resultCubeMap map[string]*enigma.HyperCube) {
	for key, value := range jsonMap {
		if key == "qHyperCube" {
			hypercube := toHypercube(value)
			resultCubeMap[key] = hypercube

		} else if !strings.HasPrefix(key, "q") {
			var customSubpath string
			if path == "" {
				customSubpath = path + "/" + key
			} else {
				customSubpath = key
			}
			switch jsonSubMap := value.(type) {
			case map[string]interface{}:
				getAllHyperCubes(customSubpath, jsonSubMap, resultCubeMap)
			}
		}
	}
}

func toHypercube(value interface{}) *enigma.HyperCube {
	subnode, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	var hypercube enigma.HyperCube
	err = json.Unmarshal(subnode, &hypercube)
	if err != nil {
		fmt.Println(err)
	}
	return &hypercube
}
