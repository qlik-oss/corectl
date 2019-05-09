package printer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/enigma-go"
)

// PrintNamedItemsList prints a list of the id and type and title of the supplied items
func PrintNamedItemsList(items []internal.NamedItem, printAsBash bool) {
	if printAsBash {
		for _, item := range items {
			fmt.Println(item.Id)
		}
	} else if internal.PrintJSON {
		internal.PrintAsJSON(items)
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Id", "Title"})
		for _, info := range items {
			writer.Append([]string{info.Id, info.Title})
		}
		writer.Render()
	}
}

// PrintNamedItemsList prints a list of the id and type and title of the supplied items
func PrintNamedItemsListWithType(items []internal.NamedItemWithType, printAsBash bool) {
	if printAsBash {
		for _, item := range items {
			fmt.Println(item.Id)
		}
	} else if internal.PrintJSON {
		internal.PrintAsJSON(items)
	} else {
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"Id", "Type", "Title"})
		for _, info := range items {
			writer.Append([]string{info.Id, info.Type, info.Title})
		}
		writer.Render()
	}
}

// PrintGenericEntityProperties prints the properties of the generic entity defined by entityID
func PrintGenericEntityProperties(state *internal.State, entityID string, entityType string, minimum bool) {
	var err error
	var properties json.RawMessage

	if minimum {
		switch entityType {
		case "object":
			genericObject, err := state.Doc.GetObject(state.Ctx, entityID)
			if err != nil {
				internal.FatalError(err)
			}
			qProps, _ := genericObject.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		case "measure":
			genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
			if err != nil {
				internal.FatalError(err)
			}
			qProps, _ := genericMeasure.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		case "dimension":
			genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
			if err != nil {
				internal.FatalError(err)
			}
			qProps, _ := genericDimension.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		}
	} else {
		switch entityType {
		case "object":
			genericObject, err := state.Doc.GetObject(state.Ctx, entityID)
			if err != nil {
				internal.FatalError(err)
			}
			properties, err = genericObject.GetPropertiesRaw(state.Ctx)
		case "measure":
			genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
			if err != nil {
				internal.FatalError(err)
			}
			properties, err = genericMeasure.GetPropertiesRaw(state.Ctx)
		case "dimension":
			genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
			if err != nil {
				internal.FatalError(err)
			}
			properties, err = genericDimension.GetPropertiesRaw(state.Ctx)
		}
	}
	if err != nil {
		internal.FatalError(err)
	}
	if len(properties) == 0 {
		internal.FatalError(fmt.Sprintf("No %s by id '%s'", entityType, entityID))
	} else {
		internal.PrintAsJSON(properties)
	}
}

// PrintGenericEntityLayout prints the layout of the object defined by objectID
func PrintGenericEntityLayout(state *internal.State, entityID string, entityType string) {
	var err error
	var layout json.RawMessage
	switch entityType {
	case "object":
		genericObject, err := state.Doc.GetObject(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		layout, err = genericObject.GetLayoutRaw(state.Ctx)
	case "measure":
		genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		layout, err = genericMeasure.GetLayoutRaw(state.Ctx)
	case "dimension":
		genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
		if err != nil {
			internal.FatalError(err)
		}
		layout, err = genericDimension.GetLayoutRaw(state.Ctx)
	}
	if err != nil {
		internal.FatalError(err)
	}
	if len(layout) == 0 {
		internal.FatalError(fmt.Sprintf("No %s by id '%s'", entityType, entityID))
	} else {
		internal.PrintAsJSON(layout)
	}
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
	if len(resultCubeMap) == 0 {
		internal.FatalError(fmt.Sprintf("Object %s contains no data\n", objectID))
	}
	for _, hypercube := range resultCubeMap {
		printHypercube(hypercube)
	}
}

func printHypercube(hypercube *enigma.HyperCube) {
	writer := tablewriter.NewWriter(os.Stdout)

	headers := []string{}
	for _, dim := range hypercube.DimensionInfo {
		headers = append(headers, dim.FallbackTitle)
	}
	for _, mes := range hypercube.MeasureInfo {
		headers = append(headers, mes.FallbackTitle)
	}

	writer.SetAutoFormatHeaders(false)
	writer.SetHeader(headers)

	for _, page := range hypercube.DataPages {
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
