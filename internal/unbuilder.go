package internal

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/qlik-oss/corectl/internal/log"
	"github.com/qlik-oss/enigma-go"
)

type (
	// UnbuildEntityProperies contains prroperties to unmarshal when looking into the meta info in entities when exporting
	UnbuildEntityProperies struct {
		QInfo struct {
			QId   string `json:"qId"`
			QType string `json:"qType"`
		} `json:"qInfo"`
		QMetaDef struct {
			Title string `json:"title"`
		} `json:"qMetaDef"`
		Visualization string `json:"visualization"`
		QProperty     *UnbuildEntityProperies
	}

	// JSONWithOrder is a container for a json struct that retains the order in which the data was originally fetched
	// Used to hold the results of parallel calls
	JSONWithOrder struct {
		JSON  json.RawMessage
		Order int
	}
)

// Regex translation:
//   \pL - Unicode group for letters, meaning all letters.
//   \d  - Digits
// Summary, match anything that is not a unicode letter, number, hyphen or underscore.
// This is to ensure that our path names are not "bonkers".
var matchAllNonAlphaNumeric = regexp.MustCompile(`[^\pL\d_-]+`)

// Unbuild exports measures, dimensions, variables, connections, objects and a config file from an app into the file system
func Unbuild(ctx context.Context, doc *enigma.Doc, global *enigma.Global, rootFolder string) {
	log.Verboseln("Exporting app to folder: " + rootFolder)
	os.MkdirAll(rootFolder, os.ModePerm)
	exportEntities(ctx, doc, rootFolder)
	exportVariables(ctx, doc, rootFolder)
	exportScript(ctx, doc, rootFolder)
	exportAppProperties(ctx, doc, rootFolder)
	exportConnections(ctx, doc, rootFolder)
	exportMainConfigFile(rootFolder)
}

func exportEntities(ctx context.Context, doc *enigma.Doc, folder string) {
	measureArray := make([]JSONWithOrder, 0)
	var measureArrayLock sync.Mutex
	dimensionArray := make([]JSONWithOrder, 0)
	var dimensionArrayLock sync.Mutex
	allInfos, _ := doc.GetAllInfos(ctx)
	waitChannel := make(chan interface{}, 10000)
	defer close(waitChannel)
	for index, item := range allInfos {
		go func(index int, item *enigma.NxInfo) {
			if dimension, _ := doc.GetDimension(ctx, item.Id); dimension != nil && dimension.Type != "" {
				props, _ := dimension.GetPropertiesRaw(ctx)
				dimensionArrayLock.Lock()
				dimensionArray = append(dimensionArray, JSONWithOrder{props, index})
				dimensionArrayLock.Unlock()
			} else if measure, _ := doc.GetMeasure(ctx, item.Id); measure != nil && measure.Type != "" {
				props, _ := measure.GetPropertiesRaw(ctx)
				measureArrayLock.Lock()
				measureArray = append(measureArray, JSONWithOrder{props, index})
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
		}(index, item)
	}
	for range allInfos {
		<-waitChannel
	}
	writeMeasures(measureArray, folder)
	writeDimensions(dimensionArray, folder)
}

func exportVariables(ctx context.Context, doc *enigma.Doc, folder string) {
	variableArray := make([]JSONWithOrder, 0)
	var variarbleArraySync sync.Mutex
	variables := ListVariables(ctx, doc)
	waitChannel := make(chan interface{}, 10000)
	defer close(waitChannel)
	for index, item := range variables {
		go func(index int, item NamedItem) {
			if variable, _ := doc.GetVariableByName(ctx, item.Title); variable != nil && variable.Handle != 0 {
				variarbleArraySync.Lock()
				props, _ := variable.GetPropertiesRaw(ctx)
				variableArray = append(variableArray, JSONWithOrder{props, index})
				variarbleArraySync.Unlock()
			}
			waitChannel <- true
		}(index, item)
	}
	for range variables {
		<-waitChannel
	}
	writeVariables(variableArray, folder)
}

func exportScript(ctx context.Context, doc *enigma.Doc, folder string) {
	script, _ := doc.GetScript(ctx)
	ioutil.WriteFile(folder+"/script.qvs", []byte(script), os.ModePerm)
	log.Verboseln("Exported script to " + folder + "/script.qvs")
}

func exportAppProperties(ctx context.Context, doc *enigma.Doc, folder string) {
	appProperties, _ := doc.GetAppProperties(ctx)
	ioutil.WriteFile(folder+"/app-properties.json", marshalOrFail(appProperties), os.ModePerm)
	log.Verboseln("Exported app properties to " + folder + "/app-properties.json")
}

func exportConnections(ctx context.Context, doc *enigma.Doc, folder string) {
	connections, _ := doc.GetConnections(ctx)
	connectionsStr := "connections:\n"
	for _, x := range connections {
		connectionsStr += "  " + x.Name + ":" + "\n"
		connectionsStr += "    type: " + x.Type + "\n"
		connectionsStr += "    connectionstring: " + x.ConnectionString + "\n"
		if x.Type != "folder" {
			connectionsStr += "    username: " + x.UserName + "\n"
			connectionsStr += "    password: " + "\n"
		}
	}

	ioutil.WriteFile(folder+"/connections.yml", []byte(connectionsStr), os.ModePerm)
	log.Verbosef("Exported %v connection(s) to %s/connections.yml", len(connections), folder)
}

func exportMainConfigFile(rootFolder string) {
	config := "script: script.qvs\n" +
		"connections: connections.yml\n" +
		"dimensions: dimensions.json\n" +
		"measures: measures.json\n" +
		"objects: objects/*.json\n" +
		"variables: variables.json\n" +
		"app-properties: app-properties.json\n"
	ioutil.WriteFile(rootFolder+"/corectl.yml", []byte(config), os.ModePerm)
}

func writeDimensions(dimensionArray []JSONWithOrder, folder string) {
	sortJSONArray(dimensionArray)
	filename := folder + "/dimensions.json"
	ioutil.WriteFile(filename, marshalOrFail(toJSONArray(dimensionArray)), os.ModePerm)
	log.Verbosef("Exported %v dimension(s) to %s/dimensions.yml", len(dimensionArray), folder)
}

func writeMeasures(measureArray []JSONWithOrder, folder string) {
	sortJSONArray(measureArray)
	filename := folder + "/measures.json"
	ioutil.WriteFile(filename, marshalOrFail(toJSONArray(measureArray)), os.ModePerm)
	log.Verbosef("Exported %v measure(s) to %s/measures.yml", len(measureArray), folder)
}

func writeVariables(variableArray []JSONWithOrder, folder string) {
	sortJSONArray(variableArray)
	filename := folder + "/variables.json"
	ioutil.WriteFile(filename, marshalOrFail(toJSONArray(variableArray)), os.ModePerm)
	log.Verbosef("Exported %v variable(s) to %s/variables.yml", len(variableArray), folder)
}

func marshalOrFail(v interface{}) json.RawMessage {
	result, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalln(err)
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

// BuildRootFolderFromTitle returns a folder name based on app title
func BuildRootFolderFromTitle(title string) string {
	title = strings.ToLower(title) + "-unbuild"
	title = matchAllNonAlphaNumeric.ReplaceAllString(title, `-`)
	return title
}

func sortJSONArray(array []JSONWithOrder) {
	sort.SliceStable(array, func(i, j int) bool {
		return array[i].Order < array[j].Order
	})
}

func toJSONArray(array []JSONWithOrder) []json.RawMessage {
	result := []json.RawMessage{}
	for _, x := range array {
		result = append(result, x.JSON)
	}
	return result
}
