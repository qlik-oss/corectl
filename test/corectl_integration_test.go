// +build integration

package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/kr/pretty"
	enigma "github.com/qlik-oss/enigma-go"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golden files")

var testAppName = "corectl_test_app.qvf"

var engineIP = flag.String("engineIP", "localhost:9076", "URL to first engine instance in docker-compose.yml i.e qix-engine-1")
var engine2IP = flag.String("engine2IP", "localhost:9176", "URL to second engine instance in docker-compose.yml i.e qix-engine-2")
var engine3IP = flag.String("engine3IP", "localhost:9276", "URL to third engine instance in docker-compose.yml i.e qix-engine-3")

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "corectl.exe"
	}

	return "corectl"
}

var binaryName = getBinaryName()

var binaryPath string

type testFile struct {
	t    *testing.T
	name string
	dir  string
}

func newGoldenFile(t *testing.T, name string) *testFile {
	return &testFile{t: t, name: name, dir: "golden"}
}

func (tf *testFile) path() string {
	tf.t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		tf.t.Fatal("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
}

func (tf *testFile) write(content string) {
	tf.t.Helper()
	err := ioutil.WriteFile(tf.path(), []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

func diff(expected, actual interface{}) []string {
	return pretty.Diff(expected, actual)
}

func (tf *testFile) load() string {
	tf.t.Helper()

	content, err := ioutil.ReadFile(tf.path())
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return string(content)
}

func setupEntities(connectToEngine string, configPath string, entityType string, entityPath string) []byte {
	cmd := exec.Command(binaryPath, []string{connectToEngine, configPath, "build", entityPath}...)
	cmd.Run()
	cmd = exec.Command(binaryPath, []string{connectToEngine, configPath, entityType, "ls", "--json"}...)
	output, _ := cmd.CombinedOutput()
	return output
}

func removeEntities(t *testing.T, connectToEngine string, configPath string, entityType string, entityId string) {
	cmd := exec.Command(binaryPath, []string{connectToEngine, configPath, entityType, "rm", entityId}...)
	output, _ := cmd.CombinedOutput()
	assert.Equal(t, "Saving app... Done\n\n", string(output))
}

func verifyNoEntities(t *testing.T, connectToEngine string, configPath string, entityType string) {
	cmd := exec.Command(binaryPath, []string{connectToEngine, configPath, entityType, "ls", "--json"}...)
	output, _ := cmd.CombinedOutput()
	assert.Equal(t, "[]\n", string(output))
}

func TestNestedObjectSupport(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	//create the nested objects
	output := setupEntities(connectToEngine, "--config=test/project2/corectl-alt.yml", "object", "--objects=test/project2/sheet.json")

	//verify that the objects are created
	var objects []*enigma.NxInfo
	err := json.Unmarshal(output, &objects)
	assert.NoError(t, err)
	assert.NotEmpty(t, objects)
	assert.NotNil(t, objects[0])
	assert.NotNil(t, objects[0].Id)
	assert.Equal(t, "a699ee97-152d-4470-9655-ae7c82d71491", objects[0].Id)
	assert.Len(t, objects, 3)

	//verify that removing the objects works
	removeEntities(t, connectToEngine, "--config=test/project2/corectl.yml", "object", objects[0].Id)

	//verify that there is no objects in the app anymore.
	verifyNoEntities(t, connectToEngine, "--config=test/project2/corectl.yml", "object")

	//remove the app as clean-up (Otherwise we might share sessions when we use that app again.)
	_ = exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "app", "rm", testAppName}...)
}

func TestConnections(t *testing.T) {
	//setup env var for project 2
	os.Setenv("CONN_TYPE", "folder")

	//create the connection
	connectToEngine := "--engine=" + *engineIP
	output := setupEntities(connectToEngine, "--config=test/project2/corectl.yml", "connection", "--connections=test/project2/connections.yml")

	//verify that the connection was created
	var connections []*enigma.Connection
	err := json.Unmarshal(output, &connections)
	assert.NoError(t, err)
	assert.NotNil(t, connections[0])
	assert.NotNil(t, connections[0].Id)

	//verify that removing the connection works
	removeEntities(t, connectToEngine, "--config=test/project2/corectl.yml", "connection", connections[0].Id)

	//verify that there is no connections in the app anymore.
	verifyNoEntities(t, connectToEngine, "--config=test/project2/corectl.yml", "connection")

	//remove the app as clean-up (Otherwise we might share sessions when we use that app again.)
	_ = exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "app", "rm", testAppName}...)
}

// TestPrecedence checks that command line flags overrides config props
func TestPrecedence(t *testing.T) {
	// Set objects, dimensions, measures and connection explicitly.
	// The information in the config should therefore be overriden.
	config := "--config=test/project5/corectl.yml"
	engine := "--engine=" + *engineIP
	flags := []string{
		"build",
		config,
		engine,
		"--objects=test/project5/o/*",
		"--dimensions=test/project5/d/*",
		"--measures=test/project5/m/*",
		"--connections=test/project5/connections.yml",
	}
	cmd := exec.Command(binaryPath, flags...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
	}

	// Cleanup
	rm_cmd := exec.Command(binaryPath, config, engine, "app", "rm", testAppName)
	defer rm_cmd.Run()

	var data []map[string]string
	entities := []string{"object", "dimension", "measure"}
	expected := []string{"my-hypercube2", "swedish-dimension", "measure-x"}
	for i, entity := range entities {
		cmd = exec.Command(binaryPath, config, engine, entity, "ls", "--json")
		output, err = cmd.Output()
		json.Unmarshal(output, &data)
		assert.Nil(t, err)
		assert.Len(t, data, 1)
		assert.Equal(t, expected[i], data[0]["qId"])
	}

	var connections []*enigma.Connection
	cmd = exec.Command(binaryPath, config, engine, "connection", "ls", "--json")
	output, err = cmd.Output()
	json.Unmarshal(output, &connections)
	assert.Nil(t, err)
	assert.Len(t, connections, 1)
	assert.Equal(t, "bogusname", connections[0].Name)
}

func setupTest(t *testing.T, tt test) func(t *testing.T, tt test) {
	if tt.initTest.setup == true {
		t.Log("\u001b[96m *** Setup *** \u001b[0m")

		args := append(tt.connectString, []string{"build"}...)
		cmd := exec.Command(binaryPath, args...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Unable to create app: %s\n", output)
		}
	}

	return func(t *testing.T, tt test) {
		if tt.initTest.teardown == true {
			t.Log("\u001b[96m *** Teardown *** \u001b[0m")

			args := append(tt.connectString, []string{"app", "rm", testAppName, "--suppress"}...)
			cmd := exec.Command(binaryPath, args...)

			t.Log("\u001b[35m Executing command:" + strings.Join(cmd.Args, " ") + "\u001b[0m")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Unable to delete app: %s\n", output)
			}
		}
	}
}

type initTest struct {
	setup    bool
	teardown bool
}

type test struct {
	name          string
	connectString []string
	command       []string
	expected      []string
	initTest
}

func TestCorectl(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	connectToEngineWithInccorectLicenseService := "--engine=" + *engine2IP
	connectToEngineABAC := "--engine=" + *engine3IP

	//Setup env vars for project 1
	os.Setenv("CORECTL_TEST_CONNECT", "corectl-test-connector")
	os.Setenv("ENGINE_URL", "localhost:9076")
	// General
	emptyConnectString := []string{}
	defaultConnectString1 := []string{"--config=test/project1/corectl.yml", connectToEngine, "--verbose=false", "--traffic=false"} // Config validation

	tests := []test{
		{"help 1", emptyConnectString, []string{""}, []string{"golden", "help-1.golden"}, initTest{false, false}},
		{"help 2", emptyConnectString, []string{"help"}, []string{"golden", "help-2.golden"}, initTest{false, false}},
		{"help 3", emptyConnectString, []string{"help", "build"}, []string{"golden", "help-3.golden"}, initTest{false, false}},

		{"project 1 - build", defaultConnectString1, []string{"build"}, []string{"Connected", "TableA <<  5 Lines fetched", "TableB <<  5 Lines fetched", "Reload finished successfully", "Saving app... Done"}, initTest{false, true}},
		{"project 1 - get tables", defaultConnectString1, []string{"tables"}, []string{"golden", "project1-tables.golden"}, initTest{true, true}},
		{"project 1 - get assoc", defaultConnectString1, []string{"assoc"}, []string{"golden", "project1-assoc.golden"}, initTest{true, true}},
		{"project 1 - get fields", defaultConnectString1, []string{"fields"}, []string{"golden", "project1-fields.golden"}, initTest{true, true}},
		{"project 1 - get field numbers", defaultConnectString1, []string{"values", "numbers"}, []string{"golden", "project1-values-numbers.golden"}, initTest{true, true}},
		{"project 1 - get meta", defaultConnectString1, []string{"meta"}, []string{"golden", "project1-meta.golden"}, initTest{true, true}},
		{"project 1 - eval", defaultConnectString1, []string{"eval", "count(numbers)", "by", "xyz"}, []string{"golden", "project1-eval-1.golden"}, initTest{true, true}},
		{"project 1 - eval", defaultConnectString1, []string{"eval", "count(numbers)"}, []string{"golden", "project1-eval-2.golden"}, initTest{true, true}},
		{"project 1 - eval", defaultConnectString1, []string{"eval", "=1+1"}, []string{"golden", "project1-eval-3.golden"}, initTest{true, true}},
		{"project 1 - eval", defaultConnectString1, []string{"eval", "1+1"}, []string{"golden", "project1-eval-4.golden"}, initTest{true, true}},
		{"project 1 - eval", defaultConnectString1, []string{"eval", "by", "numbers"}, []string{"golden", "project1-eval-5.golden"}, initTest{true, true}},
		{"project 1 - get objects", defaultConnectString1, []string{"object", "ls", "--json"}, []string{"golden", "project1-objects.golden"}, initTest{true, true}},
		{"project 1 - get object data", defaultConnectString1, []string{"object", "data", "my-hypercube"}, []string{"golden", "project1-data.golden"}, initTest{true, true}},
		{"project 1 - get object properties", defaultConnectString1, []string{"object", "properties", "my-hypercube", "--json"}, []string{"golden", "project1-properties.golden"}, initTest{true, true}},
		{"project 1 - get measures 1 as json", defaultConnectString1, []string{"measure", "ls", "--json"}, []string{"golden", "project1-measures-1.golden"}, initTest{true, true}},
		{"project 1 - get dimensions", defaultConnectString1, []string{"dimension", "ls", "--json"}, []string{"golden", "project1-dimensions.golden"}, initTest{true, true}},
		{"project 1 - get script", defaultConnectString1, []string{"script", "get"}, []string{"golden", "project1-script.golden"}, initTest{true, true}},
		{"project 1 - reload without progress", defaultConnectString1, []string{"reload", "--silent"}, []string{"golden", "project1-reload-silent.golden"}, initTest{true, true}},
		{"project 1 - reload without progress and without save", defaultConnectString1, []string{"reload", "--silent", "--no-save"}, []string{"golden", "project1-reload-silent-no-save.golden"}, initTest{true, true}},
		{"project 1 - set measures", defaultConnectString1, []string{"measure", "set", "test/project1/not-following-glob-pattern-measure.json", "--no-save"}, []string{"golden", "blank.golden"}, initTest{true, true}},
		{"project 1 - get measures 2", []string{"--config=test/project1/corectl-alt.yml", connectToEngine}, []string{"measure", "ls", "--json"}, []string{"golden", "project1-measures-2.golden"}, initTest{true, true}},
		{"project 1 - remove measures", []string{"--config=test/project1/corectl-alt.yml", connectToEngine}, []string{"measure", "rm", "measure-3", "--no-save"}, []string{"golden", "blank.golden"}, initTest{true, true}},
		{"project 1 - check measures after removal", defaultConnectString1, []string{"measure", "ls", "--json"}, []string{"golden", "project1-measures-1.golden"}, initTest{true, true}},
		{"project 1 - set script", defaultConnectString1, []string{"script", "set", "test/project1/dummy-script.qvs", "--no-save"}, []string{"golden", "blank.golden"}, initTest{true, true}},
		{"project 1 - get script after setting it", []string{"--config=test/project1/corectl-alt.yml", connectToEngine}, []string{"script", "set", "test/project1/dummy-script.qvs"}, []string{"golden", "project1-script-2.golden"}, initTest{true, true}},
		{"project 1 - traffic logging", []string{"--config=test/project1/corectl-alt.yml", connectToEngine}, []string{"script", "set", "test/project1/dummy-script.qvs", "--traffic"}, []string{"golden", "project1-traffic-log.golden"}, initTest{true, true}},

		// Verify behaviour when opening an app without data
		{"project 1 - open app without data", []string{"--config=test/project1/corectl-alt.yml", "--ttl", "0", connectToEngine}, []string{"connection", "ls", "--no-data", "--verbose"}, []string{"without data"}, initTest{true, true}},
		{"project 1 - save objects in app opened without data", []string{"--config=test/project1/corectl.yml", "--ttl", "0", connectToEngine, "--traffic=false", "--verbose=false"}, []string{"build", "--no-data"}, []string{"Saving objects in app... Done"}, initTest{false, true}},

		// Project 2 has separate connections file
		{"project 2 - build with connections", []string{connectToEngine, "-a=" + testAppName, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0"}, []string{"build", "--script=test/project2/script.qvs", "--connections=test/project2/connections.yml", "--objects=test/project2/object-*.json"}, []string{"datacsv << data 1 Lines fetched", "Reload finished successfully", "Saving app... Done"}, initTest{false, true}},
		{"project 2 - build with connections 2", []string{connectToEngine, "--config=test/project2/corectl-connectionsref.yml"}, []string{"build"}, []string{"datacsv << data 1 Lines fetched", "Reload finished successfully", "Saving app... Done"}, initTest{false, true}},
		{"project 2 - get fields ", []string{"--config=test/project2/corectl-alt.yml ", connectToEngine}, []string{"fields"}, []string{"golden", "project2-fields.golden"}, initTest{true, true}},
		{"project 2 - get data", []string{"--config=test/project2/corectl-alt.yml ", connectToEngine}, []string{"object", "data", "my-hypercube-on-commandline"}, []string{"golden", "project2-data.golden"}, initTest{true, true}},
		{"project 2 - keys", []string{"--config=test/project2/corectl-alt2.yml", connectToEngine}, []string{"keys"}, []string{"animal"}, initTest{true, true}},

		{"err project 1 - invalid-catwalk-url", defaultConnectString1, []string{"catwalk", "--catwalk-url=not-a-valid-url"}, []string{"golden", "project1-catwalk-error.golden"}, initTest{false, false}},
		{"err 2", []string{connectToEngine, "--app=nosuchapp.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0"}, []string{"eval", "count(numbers)", "by", "xyz"}, []string{"golden", "err-2.golden"}, initTest{false, false}},
		{"err 3", []string{connectToEngine, "--app=" + testAppName, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0"}, []string{"object", "data", "nosuchobject"}, []string{"golden", "err-3.golden"}, initTest{true, true}},

		{"project 1 - get status", defaultConnectString1, []string{"status"}, []string{"Connected to " + testAppName + " @ ", "The data model has 2 tables."}, initTest{true, true}},
		{"list apps json", defaultConnectString1, []string{"app", "ls", "--json"}, []string{"\"id\": \"/apps/" + testAppName + "\","}, initTest{true, true}},
		{"err 1", []string{"--app=bogus", "--engine=localhost:9999"}, []string{"fields"}, []string{"Please check the --engine parameter or your config file", "Error details:  dial tcp"}, initTest{false, false}},

		// trying to connect to an engine that has JWT authorization activated without a JWT Header
		{"err jwt", []string{connectToEngine}, []string{"app", "ls"}, []string{"Error details:  401 from ws server: websocket: bad handshake"}, initTest{false, false}},
		{"err no license", []string{connectToEngineWithInccorectLicenseService}, []string{"app", "ls"}, []string{"Failed to connect to engine with error message:  SESSION_ERROR_NO_LICENSE"}, initTest{false, false}},

		// Verifying corectl against an engine running with ABAC enabled
		{"project 4 - get status", []string{"--config=test/project4/corectl.yml ", connectToEngineABAC}, []string{"status"}, []string{"Connected to " + testAppName + " @ ", "The data model has 1 table."}, initTest{true, true}},
		{"project 4 - list apps", []string{"--config=test/project4/corectl.yml ", connectToEngineABAC}, []string{"app", "ls", "--json"}, []string{"\"title\": \"" + testAppName + "\","}, initTest{true, true}},
		{"project 4 - get meta", []string{"--config=test/project4/corectl.yml ", connectToEngineABAC}, []string{"meta"}, []string{"golden", "project4-meta.golden"}, initTest{true, true}},

		// Verifying config validation
		{"err invalid 1", []string{"--config=test/project2/corectl-invalid.yml ", connectToEngine}, []string{"build"}, []string{"apps", "header", "object", "measure", "verbos", "trafic", "connection", "dimension"}, initTest{false, false}},
		{"err invalid 2", []string{"--config=test/project2/corectl-invalid2.yml ", connectToEngine}, []string{"build"}, []string{"'header': did you mean 'headers'?", "test/project2/corectl-invalid2.yml"}, initTest{false, false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup and teardown
			teardownTest := setupTest(t, tt)
			defer teardownTest(t, tt)

			args := append(tt.connectString, tt.command...)
			cmd := exec.Command(binaryPath, args...)

			t.Log("\u001b[35m Executing command:" + strings.Join(cmd.Args, " ") + "\u001b[0m")
			output, err := cmd.CombinedOutput()

			if strings.HasPrefix(tt.name, "err") {
				if err == nil {
					t.Fatalf("%s\nexpected (err == nil) to be %v, but got %v. err: %v", output, false, err == nil, err)
				}
			} else if err != nil {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, false, err != nil, err)
			}
			actual := string(output)

			if tt.expected[0] == "golden" {
				golden := newGoldenFile(t, tt.expected[1])

				if *update {
					golden.write(actual)
				}
				expected := golden.load()

				if !reflect.DeepEqual(expected, actual) {
					t.Fatalf("diff: %v", diff(expected, actual))
				}
			} else {
				for _, sub := range tt.expected {
					if !strings.Contains(actual, sub) {
						t.Fatalf("Output did not contain substring '%v'\nReceived:\n%v", sub, actual)
					}
				}
			}
		})
	}
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}

	abs, err := filepath.Abs(binaryName)
	if err != nil {
		fmt.Printf("could not get abs path for %s: %v", binaryName, err)
		os.Exit(1)
	}

	binaryPath = abs

	if err := exec.Command("go", "build", "-o", binaryName, "-v").Run(); err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
