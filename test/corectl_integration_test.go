// +build integration

package test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"reflect"

	"github.com/kr/pretty"
)

var update = flag.Bool("update", false, "update golden files")

var engineIP = flag.String("engineIP", "localhost:9076", "dir of package containing embedded files")

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

func TestCorectlGolden(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	tests := []struct {
		name   string
		args   []string
		golden string
	}{
		{"help 1", []string{""}, "help-1.golden"},
		{"help 2", []string{"help"}, "help-2.golden"},
		{"help 3", []string{"help", "build"}, "help-3.golden"},
		{"project 1 - build", []string{"--config=test/project1/corectl.yml", connectToEngine, "build"}, "project1-build.golden"},
		{"project 1 - get tables", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "tables"}, "project1-tables.golden"},
		{"project 1 - get assoc", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "assoc"}, "project1-assoc.golden"},
		{"project 1 - get fields", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "fields"}, "project1-fields.golden"},
		{"project 1 - get field numbers", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "field", "numbers"}, "project1-field-numbers.golden"},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "count(numbers)", "by", "xyz"}, "project1-eval-1.golden"},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "count(numbers)"}, "project1-eval-2.golden"},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "=1+1"}, "project1-eval-3.golden"},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "1+1"}, "project1-eval-4.golden"},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "by", "numbers"}, "project1-eval-5.golden"},
		{"project 1 - get objects", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "objects"}, "project1-objects.golden"},
		{"project 1 - get object data", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "object", "data", "my-hypercube"}, "project1-data.golden"},
		{"project 1 - get object properties", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "object", "properties", "my-hypercube"}, "project1-properties.golden"},
		{"project 1 - get object", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "object", "my-hypercube"}, "project1-properties.golden"},
		{"project 1 - get measures 1", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures"}, "project1-measures-1.golden"},
		{"project 1 - get dimensions", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "dimensions"}, "project1-dimensions.golden"},
		{"project 1 - get script", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "script"}, "project1-script.golden"},
		{"project 1 - get status", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "status"}, "project1-status.golden"},
		{"project 1 - reload without progress", []string{"--config=test/project1/corectl.yml", connectToEngine, "reload", "--silent"}, "project1-reload-silent.golden"},
		{"project 1 - reload without progress and without save", []string{"--config=test/project1/corectl.yml", connectToEngine, "reload", "--silent", "--noSave"}, "project1-reload-silent-nosave.golden"},
		{"project 1 - set measures", []string{"--config=test/project1/corectl.yml ", connectToEngine, "set", "measures", "--measures=test/project1/not-following-glob-pattern-measure.json", "--noSave"}, "blank.golden"},
		{"project 1 - get measures 2", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures"}, "project1-measures-2.golden"},
		{"project 1 - remove measures", []string{"--config=test/project1/corectl.yml ", connectToEngine, "remove", "measures", "measure-3", "--noSave"}, "blank.golden"},
		{"project 1 - check measures after removal", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "measures"}, "project1-measures-1.golden"},
		{"project 1 - set script", []string{"--config=test/project1/corectl.yml ", connectToEngine, "set", "script", "--script=test/project1/dummy-script.qvs", "--noSave"}, "blank.golden"},
		{"project 1 - get script after setting it", []string{"--config=test/project1/corectl.yml ", connectToEngine, "get", "script"}, "project1-script-2.golden"},

		// Project 2 has separate connections file
		{"project 2 - build with connections", []string{connectToEngine, "-a=project2.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "build", "--script=test/project2/script.qvs", "--connections=test/project2/connections.yml", "--objects=test/project2/object-*.json"}, "project2-build.golden"},
		{"project 2 - get fields ", []string{"--config=test/project2/corectl.yml ", connectToEngine, "get", "fields"}, "project2-fields.golden"},
		{"project 2 - get data", []string{"--config=test/project2/corectl.yml ", connectToEngine, "get", "object", "data", "my-hypercube-on-commandline"}, "project2-data.golden"},

		{"project 3 - build ", []string{"--config=test/project3/corectl.yml ", connectToEngine, "build"}, "project3-build.golden"},
		{"project 3 - get fields", []string{"--config=test/project3/corectl.yml ", connectToEngine, "get", "fields"}, "project3-fields.golden"},
		{"err 2", []string{connectToEngine, "--app=nosuchapp.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "eval", "count(numbers)", "by", "xyz"}, "err-2.golden"},
		{"err 3", []string{connectToEngine, "--app=project1.qvf", "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "get", "object", "data", "nosuchobject"}, "err-3.golden"},
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

			golden := newGoldenFile(t, tt.golden)

			if *update {
				golden.write(actual)
			}
			expected := golden.load()

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		})
	}
}

func TestCorectlContains(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{"list apps", []string{connectToEngine, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "get", "apps"}, []string{"Id", "Name", "Last-Reloaded", "ReadOnly", "Title", "project2.qvf", "project1.qvf"}},
		{"list apps json", []string{connectToEngine, "--headers=authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmb2xrZSJ9.MD_revuZ8lCEa6bb-qtfYaHdxBiRMUkuH86c4kd1yC0", "get", "apps", "--json"}, []string{"\"id\": \"/apps/project2.qvf\","}},
		{"err 1", []string{"--engine=localhost:9999", "get", "fields"}, []string{"Please check the --engine parameter or your config file", "Error details:  dial tcp"}},
		// trying to connect to an engine that has JWT authorization activated without a JWT Header
		{"err jwt", []string{connectToEngine, "get", "apps"}, []string{"Error details:  401 from ws server: websocket: bad handshake"}},
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

			for _, sub := range tt.contains {
				if !strings.Contains(actual, sub) {
					t.Fatalf("Output did not contain substring %v", sub)
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
