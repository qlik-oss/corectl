package internal

import (
	"context"
	"encoding/json"
	"github.com/qlik-oss/enigma-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type (
	// UnbuildEntityProperies contains prroperties to unmarshal when looking into the meta info in entities when exporting
	UnbuildEntityProperies struct {
		QInfo struct {
			QId   string `json:"qId"`
			QType string `json:"qType"`
		} `jsom: qInfo`
		QMetaDef struct {
			Title string `json:"title"`
		} `jsom: qMetaDef`
		Visualization string `jsom: visualization`
		QProperty     *UnbuildEntityProperies
	}
)

var matchAllNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Unbuild exports measures, dimensions, variables, connections, objects and a config file from an app into the file system
func Unbuild(ctx context.Context, doc *enigma.Doc, global *enigma.Global, rootFolder string) {
	os.MkdirAll(rootFolder, os.ModePerm)
	exportEntities(ctx, doc, rootFolder)
	exportVariables(ctx, doc, rootFolder)
	exportScript(ctx, doc, rootFolder)
	connectionsStr := exportConnections(ctx, doc)
	exportMainConfigFile(connectionsStr, rootFolder)
}

func exportEntities(ctx context.Context, doc *enigma.Doc, folder string) {
	measureArray := make([]json.RawMessage, 0)
	var measureArrayLock sync.Mutex
	dimensionArray := make([]json.RawMessage, 0)
	var dimensionArrayLock sync.Mutex
	allInfos, _ := doc.GetAllInfos(ctx)
	waitChannel := make(chan interface{}, 10000)
	defer close(waitChannel)
	for _, item := range allInfos {
		go func(item *enigma.NxInfo) {
			if dimension, _ := doc.GetDimension(ctx, item.Id); dimension != nil && dimension.Type != "" {
				props, _ := dimension.GetPropertiesRaw(ctx)
				dimensionArrayLock.Lock()
				dimensionArray = append(dimensionArray, props)
				dimensionArrayLock.Unlock()
			} else if measure, _ := doc.GetMeasure(ctx, item.Id); measure != nil && measure.Type != "" {
				props, _ := measure.GetPropertiesRaw(ctx)
				measureArrayLock.Lock()
				measureArray = append(measureArray, props)
				measureArrayLock.Unlock()
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
					propsWithTitle := &UnbuildEntityProperies{}
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
					filename := buildEntityFilename(folder+"/objects", qType, viz, title)
					os.MkdirAll(filepath.Dir(filename), os.ModePerm)
					ioutil.WriteFile(filename, marshalOrFail(rawProps), os.ModePerm)
				}
			}
			waitChannel <- true
		}(item)
	}
	for range allInfos {
		<-waitChannel
	}
	writeMeasures(measureArrayLock, measureArray, folder)
	writeDimensions(dimensionArrayLock, dimensionArray, folder)
}

func exportVariables(ctx context.Context, doc *enigma.Doc, folder string) {
	variableArray := make([]json.RawMessage, 0)
	var variarbleArraySync sync.Mutex
	variables := ListVariables(ctx, doc)
	waitChannel := make(chan interface{}, 10000)
	defer close(waitChannel)
	for _, item := range variables {
		go func(item NamedItem) {
			if variable, _ := doc.GetVariableByName(ctx, item.Title); variable != nil && variable.Handle != 0 {
				variarbleArraySync.Lock()
				props, _ := variable.GetPropertiesRaw(ctx)
				variableArray = append(variableArray, props)
				variarbleArraySync.Unlock()
			} else if dimension, _ := doc.GetDimension(ctx, item.Id); dimension != nil && dimension.Type != "" {
			}
			waitChannel <- true
		}(item)
	}
	for range variables {
		<-waitChannel
	}
	writeVariables(variarbleArraySync, variableArray, folder)
}

func exportScript(ctx context.Context, doc *enigma.Doc, folder string) {
	script, _ := doc.GetScript(ctx)
	ioutil.WriteFile(folder+"/script.qvs", []byte(script), os.ModePerm)
}

func exportConnections(ctx context.Context, doc *enigma.Doc) string {
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
	return connectionsStr
}

func exportMainConfigFile(connectionsStr string, rootFolder string) {
	config := `script: script.qvs
		` + connectionsStr + `dimensions: dimensions.json
		measures: measures.json
		objects: objects/*.json
		variables: variables/*.json
		`
	ioutil.WriteFile(rootFolder+"/corectl.yml", []byte(config), os.ModePerm)
}

func writeMeasures(measureArrayLock sync.Mutex, measureArray []json.RawMessage, folder string) {
	measureArrayLock.Lock()
	filename := folder + "/measures.json"
	ioutil.WriteFile(filename, marshalOrFail(measureArray), os.ModePerm)
	measureArrayLock.Unlock()
}

func writeVariables(variableArrayLock sync.Mutex, variableArray []json.RawMessage, folder string) {
	variableArrayLock.Lock()
	filename := folder + "/variables.json"
	ioutil.WriteFile(filename, marshalOrFail(variableArray), os.ModePerm)
	variableArrayLock.Unlock()
}

func writeDimensions(dimensionArrayLock sync.Mutex, dimensionArray []json.RawMessage, folder string) {
	dimensionArrayLock.Lock()
	filename := folder + "/dimensions.json"
	ioutil.WriteFile(filename, marshalOrFail(dimensionArray), os.ModePerm)
	dimensionArrayLock.Unlock()
}

func marshalOrFail(v interface{}) json.RawMessage {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		FatalError(err)
	}
	return json.RawMessage(result)
}

func buildEntityFilename(folder, qType, viz, title string) string {
	qType = strings.Replace(qType, "/", "-", -1)
	viz = strings.Replace(viz, "/", "-", -1)
	title = strings.Replace(title, "/", "-", -1)
	filename := qType + "-" + viz + "-" + title
	filename = strings.ToLower(filename)
	filename = matchAllNonAlphaNumeric.ReplaceAllString(filename, `-`)
	return folder + "/" + filename + ".json"
}
