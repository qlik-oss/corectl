package printer

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qlik-oss/corectl/internal"
	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
	"github.com/spf13/viper"
)

// PrintNamedItemsList prints a list of the id and type and title of the supplied items
func PrintNamedItemsList(items []internal.NamedItem, printAsBash bool, printTitle bool) {
	switch mode {
	case jsonMode:
		log.PrintAsJSON(items)
	case bashMode:
		fallthrough
	case quietMode:
		if printTitle {
			for _, item := range items {
				log.Quietln(item.Title)
			}
		} else {
			for _, item := range items {
				log.Quietln(item.ID)
			}
		}
	default:
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"ID", "Title"})
		for _, info := range items {
			writer.Append([]string{info.ID, info.Title})
		}
		writer.Render()
	}
}

// PrintNamedItemsListWithType prints a list of the id and type and title of the supplied items
func PrintNamedItemsListWithType(items []internal.NamedItemWithType, printAsBash bool) {
	switch mode {
	case jsonMode:
		log.PrintAsJSON(items)
	case bashMode:
		fallthrough
	case quietMode:
		for _, item := range items {
			log.Quietln(item.ID)
		}
	default:
		writer := tablewriter.NewWriter(os.Stdout)
		writer.SetAutoFormatHeaders(false)
		writer.SetHeader([]string{"ID", "Type", "Title"})
		for _, info := range items {
			writer.Append([]string{info.ID, info.Type, info.Title})
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
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			full := viper.GetBool("full")
			if full {
				qPropsFull, _ := genericObject.GetFullPropertyTree(state.Ctx)
				properties, _ = json.Marshal(qPropsFull)
			} else {
				qProps, _ := genericObject.GetProperties(state.Ctx)
				properties, _ = json.Marshal(qProps)
			}
		case "measure":
			genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			qProps, _ := genericMeasure.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		case "dimension":
			genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			qProps, _ := genericDimension.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		case "variable":
			//In this case we need name not ID
			genericVariable, err := state.Doc.GetVariableByName(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			qProps, _ := genericVariable.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		case "bookmark":
			genericBookmark, err := state.Doc.GetBookmark(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			qProps, _ := genericBookmark.GetProperties(state.Ctx)
			properties, _ = json.Marshal(qProps)
		}
	} else {
		switch entityType {
		case "object":
			genericObject, err := state.Doc.GetObject(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}

			full := viper.GetBool("full")
			if full {
				properties, err = genericObject.GetFullPropertyTreeRaw(state.Ctx)
			} else {
				properties, err = genericObject.GetPropertiesRaw(state.Ctx)
			}
		case "measure":
			genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			properties, err = genericMeasure.GetPropertiesRaw(state.Ctx)
		case "dimension":
			genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			properties, err = genericDimension.GetPropertiesRaw(state.Ctx)
		case "variable":
			//In this case we need name not ID
			genericVariable, err := state.Doc.GetVariableByName(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			properties, err = genericVariable.GetPropertiesRaw(state.Ctx)
		case "bookmark":
			genericBookmark, err := state.Doc.GetBookmark(state.Ctx, entityID)
			if err != nil {
				log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
			}
			properties, err = genericBookmark.GetPropertiesRaw(state.Ctx)
		}
	}
	if err != nil {
		log.Fatalf("could not print properties of %s by ID '%s': %s\n", entityType, entityID, err)
	}
	if len(properties) == 0 {
		log.Fatalf("no %s by ID '%s'\n", entityType, entityID)
	} else {
		log.PrintAsJSON(properties)
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
			log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
		}
		layout, err = genericObject.GetLayoutRaw(state.Ctx)
	case "measure":
		genericMeasure, err := state.Doc.GetMeasure(state.Ctx, entityID)
		if err != nil {
			log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
		}
		layout, err = genericMeasure.GetLayoutRaw(state.Ctx)
	case "dimension":
		genericDimension, err := state.Doc.GetDimension(state.Ctx, entityID)
		if err != nil {
			log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
		}
		layout, err = genericDimension.GetLayoutRaw(state.Ctx)
	case "variable":
		// In this case we need name not ID
		genericVariable, err := state.Doc.GetVariableByName(state.Ctx, entityID)
		if err != nil {
			log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
		}
		layout, err = genericVariable.GetLayoutRaw(state.Ctx)
	case "bookmark":
		genericBookmark, err := state.Doc.GetBookmark(state.Ctx, entityID)
		if err != nil {
			log.Fatalf("could not retrieve %s by ID '%s': %s\n", entityType, entityID, err)
		}
		layout, err = genericBookmark.GetLayoutRaw(state.Ctx)
	}
	if err != nil {
		log.Fatalf("could not get layout of %s by ID '%s': %s\n", entityType, entityID, err)
	}
	if len(layout) == 0 {
		log.Fatalf("no %s by ID '%s'\n", entityType, entityID)
	} else {
		log.PrintAsJSON(layout)
	}
}

// EvalObject evalutes the data of the object identified by objectID
func EvalObject(ctx context.Context, doc *enigma.Doc, objectID string) {
	object, err := doc.GetObject(ctx, objectID)
	if err != nil {
		log.Fatalf("could not retrieve object by ID '%s': %s\n", objectID, err)
	}

	layout, err := object.GetLayoutRaw(ctx)
	if err != nil {
		log.Fatalf("could not get layout of object by ID '%s': %s\n", objectID, err)
	}

	layoutMap := make(map[string]interface{})
	err = json.Unmarshal(layout, &layoutMap)
	if err != nil {
		log.Errorln(err)
	}
	resultCubeMap := make(map[string]*enigma.HyperCube)
	getAllHyperCubes("", layoutMap, resultCubeMap)
	if len(resultCubeMap) == 0 {
		log.Fatalf("object %s contains no data\n", objectID)
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
		log.Errorln(err)
	}
	var hypercube enigma.HyperCube
	err = json.Unmarshal(subnode, &hypercube)
	if err != nil {
		log.Errorln(err)
	}
	return &hypercube
}
