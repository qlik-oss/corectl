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

var engineIP = flag.String("engineIP", "localhost:9076", "dir of package containing embedded files")
var engine2IP = flag.String("engine2IP", "localhost:9176", "dir of package containing embedded files")

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "corectl.exe"
	} else {
		return "corectl"
	}
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

func TestConnections(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	cmd := exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "build", "--connections=test/project2/connections.yml"}...)
	cmd.Run()
	cmd = exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "get", "connections", "--json"}...)
	output, _ := cmd.CombinedOutput()

	//verify that the connection was created
	var connections []*enigma.Connection
	err := json.Unmarshal(output, &connections)
	assert.NoError(t, err)
	assert.NotNil(t, connections[0])
	assert.NotNil(t, connections[0].Id)

	//verify that removing the connection works
	cmd = exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "remove", "connection", connections[0].Id}...)
	output, _ = cmd.CombinedOutput()
	assert.Equal(t, "Saving...Done\n\n", string(output))

	//verify that there is no connections in the app anymore.
	cmd = exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "get", "connections", "--json"}...)
	output, _ = cmd.CombinedOutput()
	assert.Equal(t, "[]\n", string(output))

	//remove the app as clean-up (Otherwise we might share sessions when we use that app again.)
	_ = exec.Command(binaryPath, []string{connectToEngine, "--config=test/project2/corectl.yml", "remove", "app", "project1.qvf"}...)
}

func TestCorectl(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	connectToEngineWithInccorectLicenseService := "--engine=" + *engine2IP
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{"help 1", []string{""}, []string{"golden", "help-1.golden"}},
		{"help 2", []string{"help"}, []string{"golden", "help-2.golden"}},
		{"help 3", []string{"help", "build"}, []string{"golden", "help-3.golden"}},
		{"project 1 - build", []string{"--config=test/project1/corectl.yml", connectToEngine, "build"}, []string{"Connected", "TableA <<  5 Lines fetched", "TableB <<  5 Lines fetched", "Reload finished successfully", "Saving...Done"}},
		{"project 1 - get tables", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "tables"}, []string{"golden", "project1-tables.golden"}},
		{"project 1 - get assoc", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "assoc"}, []string{"golden", "project1-assoc.golden"}},
		{"project 1 - get fields", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "fields"}, []string{"golden", "project1-fields.golden"}},
		{"project 1 - get field numbers", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "field", "numbers"}, []string{"golden", "project1-field-numbers.golden"}},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "count(numbers)", "by", "xyz"}, []string{"golden", "project1-eval-1.golden"}},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "count(numbers)"}, []string{"golden", "project1-eval-2.golden"}},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "=1+1"}, []string{"golden", "project1-eval-3.golden"}},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "1+1"}, []string{"golden", "project1-eval-4.golden"}},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "by", "numbers"}, []string{"golden", "project1-eval-5.golden"}},
		{"project 1 - get objects", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "objects"}, []string{"golden", "project1-objects.golden"}},
		{"project 1 - get object data", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "object", "data", "my-hypercube"}, []string{"golden", "project1-data.golden"}},
		{"project 1 - get object properties", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "object", "properties", "my-hypercube"}, []string{"golden", "project1-properties.golden"}},
		{"project 1 - get object", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "object", "my-hypercube"}, []string{"golden", "project1-properties.golden"}},
		{"project 1 - get measures 1", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures"}, []string{"golden", "project1-measures-1.golden"}},
		{"project 1 - get measures 1 as json", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures", "--json"}, []string{"golden", "project1-measures-1-json.golden"}},
		{"project 1 - get dimensions", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "dimensions"}, []string{"golden", "project1-dimensions.golden"}},
		{"project 1 - get script", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "script"}, []string{"golden", "project1-script.golden"}},
		{"project 1 - reload without progress", []string{"--config=test/project1/corectl.yml", connectToEngine, "reload", "--silent"}, []string{"golden", "project1-reload-silent.golden"}},
		{"project 1 - reload without progress and without save", []string{"--config=test/project1/corectl.yml", connectToEngine, "reload", "--silent", "--no-save"}, []string{"golden", "project1-reload-silent-no-save.golden"}},
		{"project 1 - set measures", []string{"--config=test/project1/corectl.yml ", connectToEngine, "set", "measures", "test/project1/not-following-glob-pattern-measure.json", "--no-save"}, []string{"golden", "blank.golden"}},
		{"project 1 - get measures 2", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures"}, []string{"golden", "project1-measures-2.golden"}},
		{"project 1 - remove measures", []string{"--config=test/project1/corectl.yml ", connectToEngine, "remove", "measures", "measure-3", "--no-save"}, []string{"golden", "blank.golden"}},
		{"project 1 - check measures after removal", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures"}, []string{"golden", "project1-measures-1.golden"}},
		{"project 1 - set script", []string{"--config=test/project1/corectl.yml ", connectToEngine, "set", "script", "test/project1/dummy-script.qvs", "--no-save"}, []string{"golden", "blank.golden"}},
		{"project 1 - get script after setting it", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "script"}, []string{"golden", "project1-script-2.golden"}},
		{"project 1 - traffic logging", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "script", "--traffic"}, []string{"golden", "project1-traffic-log.golden"}},

		// Project 2 has separate connections file
		{"project 2 - build with connections", []string{connectToEngine, "-a=project2.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "build", "--script=test/project2/script.qvs", "--connections=test/project2/connections.yml", "--objects=test/project2/object-*.json"}, []string{"datacsv << data 1 Lines fetched", "Reload finished successfully", "Saving...Done"}},
		{"project 2 - get fields ", []string{"--config=test/project2/corectl.yml ", connectToEngine, "get", "fields"}, []string{"golden", "project2-fields.golden"}},
		{"project 2 - get data", []string{"--config=test/project2/corectl.yml ", connectToEngine, "get", "object", "data", "my-hypercube-on-commandline"}, []string{"golden", "project2-data.golden"}},

		{"project 3 - build ", []string{"--config=test/project3/corectl.yml ", connectToEngine, "build"}, []string{"No app specified, using session app.", "datacsv << data 1 Lines fetched", "Reload finished successfully"}},
		{"project 3 - get fields", []string{"--config=test/project3/corectl.yml ", connectToEngine, "get", "fields"}, []string{"golden", "project3-fields.golden"}},
		{"err project 1 - invalid-catwalk-url", []string{"--config=test/project1/corectl.yml", connectToEngine, "catwalk", "--catwalk-url=not-a-valid-url"}, []string{"golden", "project1-catwalk-error.golden"}},
		{"err 2", []string{connectToEngine, "--app=nosuchapp.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "eval", "count(numbers)", "by", "xyz"}, []string{"golden", "err-2.golden"}},
		{"err 3", []string{connectToEngine, "--app=project1.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "get", "object", "data", "nosuchobject"}, []string{"golden", "err-3.golden"}},

		{"project 1 - get status", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "status"}, []string{"Connected to project1.qvf @ ", "The data model has 2 tables."}},
		{"list apps", []string{connectToEngine, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "get", "apps"}, []string{"Id", "Name", "Last-Reloaded", "ReadOnly", "Title", "project2.qvf", "project1.qvf"}},
		{"list apps json", []string{connectToEngine, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "get", "apps", "--json"}, []string{"\"id\": \"/apps/project2.qvf\","}},
		{"err 1", []string{"--engine=localhost:9999", "get", "fields"}, []string{"Please check the --engine parameter or your config file", "Error details:  dial tcp"}},
		// trying to connect to an engine that has JWT authorization activated without a JWT Header
		{"err jwt", []string{connectToEngine, "get", "apps"}, []string{"Error details:  401 from ws server: websocket: bad handshake"}},
		{"err no license", []string{connectToEngineWithInccorectLicenseService, "get", "apps"}, []string{"Failed to connect to engine with error message:  SESSION_ERROR_NO_LICENSE"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
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
						t.Fatalf("Output did not contain substring %v", sub)
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
