package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/qlik-oss/enigma-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type (
	Propz struct {
		QInfo struct {
			QId   string `json:"qId"`
			QType string `json:"qType"`
		} `jsom: qInfo`
		QMetaDef struct {
			Title string `json:"title"`
		} `jsom: qMetaDef`
		Visualization string `jsom: visualization`
		QProperty     *Propz
	}
)

func Unbuild(ctx context.Context, doc *enigma.Doc, global *enigma.Global) {
	rootFolder := "./unbuild-output"
	os.MkdirAll(rootFolder, os.ModePerm)

	exportEntities(ctx, doc, rootFolder)

	exportScript(ctx, doc, rootFolder)

	connections, _ := doc.GetConnections(ctx)
	connectionsStr := "connections:\n"
	for _, x := range connections {
		connectionsStr += "  " + x.Name + ": " + "\n"
		connectionsStr += "    type: " + x.Type + "\n"
		connectionsStr += "    connectionstring: " + x.ConnectionString + "\n"
		if x.Type != "folder" {
			connectionsStr += "    username: " + x.UserName + "\n"
			connectionsStr += "    password: " + "\n"
		}
	}

	config := `script: script.qvs
` + connectionsStr +
		`dimensions: dimensions.json
measures: measures.json
objects: objects/*.json
`
	writeFile(rootFolder+"/corectl.yml", config)
}

func exportEntities(ctx context.Context, doc *enigma.Doc, folder string) {
	measureArray := make([]json.RawMessage, 0)
	var measureArrayLock sync.Mutex
	dimensionArray := make([]json.RawMessage, 0)
	var dimensionArrayLock sync.Mutex
	variables, _ := doc.GetVariables(ctx, &enigma.VariableListDef{})
	fmt.Println(variables)
	allInfos, _ := doc.GetAllInfos(ctx)
	waitChannel := make(chan *NamedItemWithType, 10000)
	defer close(waitChannel)
	for _, item := range allInfos {
		go func(item *enigma.NxInfo) {
			if bookmark, _ := doc.GetBookmark(ctx, item.Id); bookmark != nil && bookmark.Type != "" {

			} else if dimension, _ := doc.GetDimension(ctx, item.Id); dimension != nil && dimension.Type != "" {
				fmt.Println("dimension")
				props, _ := dimension.GetPropertiesRaw(ctx)
				dimensionArrayLock.Lock()
				dimensionArray = append(dimensionArray, props)
				dimensionArrayLock.Unlock()
			} else if measure, _ := doc.GetMeasure(ctx, item.Id); measure != nil && measure.Type != "" {
				fmt.Println("measure")
				props, _ := measure.GetPropertiesRaw(ctx)
				measureArrayLock.Lock()
				measureArray = append(measureArray, props)
				measureArrayLock.Unlock()
			} else if variable, _ := doc.GetVariableById(ctx, item.Id); variable != nil && variable.Type != "" {
				fmt.Println("variable")
			} else if object, _ := doc.GetObject(ctx, item.Id); object != nil && object.Type != "" {
				parent, _ := object.GetParent(ctx)
				children, _ := object.GetChildInfos(ctx)
				if parent.Handle == 0 {
					var rawProps json.RawMessage
					if len(children) > 0 {
						rawProps, _ = object.GetFullPropertyTreeRaw(ctx)
					} else {
						rawProps, _ = object.GetPropertiesRaw(ctx)
					}
					propsWithTitle := &Propz{}
					json.Unmarshal(rawProps, propsWithTitle)
					if propsWithTitle.QProperty != nil {
						propsWithTitle = propsWithTitle.QProperty
					}
					title := propsWithTitle.QMetaDef.Title
					if title == "" {
						title = propsWithTitle.QInfo.QId
					}
					qType := propsWithTitle.QInfo.QType
					viz := propsWithTitle.Visualization
					filename := buildName(folder+"/objects", qType, viz, title)
					poor(filename, rawProps)
				} else {
					fmt.Println("Child object")
				}

			} else {
				fmt.Println("Unknown type")
			}
			waitChannel <- &NamedItemWithType{}
		}(item)
	}
	//Put all responses into a map by their Id
	for range allInfos {
		<-waitChannel
	}
	writeMeasures(measureArrayLock, folder, measureArray)
	writeDimensions(dimensionArrayLock, folder, dimensionArray)
}

//func exportVariables(ctx context.Context, doc *enigma.Doc, folder string) {
//	variableArray := make([]json.RawMessage, 0)
//	var variarbleArraySync sync.Mutex
//	variables, _ := doc.GetVariables(ctx, &enigma.VariableListDef{ Type:"variable", Data: json.RawMessage(`{"name":"/qMetaDef/name"}`),})
//	fmt.Println(variables)
//	waitChannel := make(chan *NamedItemWithType, 10000)
//	defer close(waitChannel)
//	for _, item := range variables {
//		go func(item *enigma.NxInfo) {
//			result := []NamedItem{}
//
//			if variable, _ := doc.GetVariable(ctx, item.Id) {
//				variarbleArraySync.Lock()
//				props, _ := variable.GetPropertiesRaw(ctx)
//				variableArray = append(variableArray, props)
//				variarbleArraySync.Unlock()
//			} else if dimension, _ := doc.GetDimension(ctx, item.Id); dimension != nil && dimension.Type != "" {
//			}
//			waitChannel <- &NamedItemWithType{}
//		}(item)
//	}
//	//Put all responses into a map by their Id
//	for range allInfos {
//		<-waitChannel
//	}
//	writeMeasures(measureArrayLock, folder, measureArray)
//	writeDimensions(variarbleArraySync, folder, variableArray)
//}
//}

func exportScript(ctx context.Context, doc *enigma.Doc, folder string) {
	script, _ := doc.GetScript(ctx)
	writeFile(folder+"/script.qvs", script)
}

func writeFile(filename string, script string) {
	ioutil.WriteFile(filename, []byte(script), os.ModePerm)
}

var matchAllNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func buildName(folder, qType, viz, title string) string {
	qType = strings.Replace(qType, "/", "-", -1)
	viz = strings.Replace(viz, "/", "-", -1)
	title = strings.Replace(title, "/", "-", -1)
	filename := qType + "-" + viz + "-" + title
	filename = strings.ToLower(filename)
	filename = matchAllNonAlphaNumeric.ReplaceAllString(filename, `-`)
	fmt.Println("------------", filename)
	return folder + "/" + filename + ".json"
}

func writeMeasures(measureArrayLock sync.Mutex, folder string, measureArray []json.RawMessage) {
	measureArrayLock.Lock()
	filename := folder + "/measures.json"
	ioutil.WriteFile(filename, marshal(measureArray), os.ModePerm)
	measureArrayLock.Unlock()
}
func writeDimensions(dimensionArrayLock sync.Mutex, folder string, dimensionArray []json.RawMessage) {
	dimensionArrayLock.Lock()
	filename := folder + "/dimensions.json"
	ioutil.WriteFile(filename, marshal(dimensionArray), os.ModePerm)
	dimensionArrayLock.Unlock()
}

func marshal(v interface{}) json.RawMessage {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		FatalError(err)
	}

	return json.RawMessage(result)
}
func poor(filename string, data json.RawMessage) {
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	var formatted bytes.Buffer
	json.Indent(&formatted, data, "", "  ")
	ioutil.WriteFile(filename, []byte(formatted.String()), os.ModePerm)
}
