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
		{"help 1", []string{"--config=test/project1/corectl.yml", "--engine=localhost:9999", ""}, "help-1.golden"},
		{"help 2", []string{"--config=test/project1/corectl.yml", "--engine=localhost:9999", "help"}, "help-2.golden"},
		{"help 3", []string{"--config=test/project1/corectl.yml", "--engine=localhost:9999", "help", "reload"}, "help-3.golden"},
		{"project 1 - reload", []string{"--config=test/project1/corectl.yml", connectToEngine, "reload"}, "project1-reload.golden"},
		{"project 1 - tables", []string{"--config=test/project1/corectl.yml ", connectToEngine, "tables"}, "project1-tables.golden"},
		{"project 1 - assoc", []string{"--config=test/project1/corectl.yml ", connectToEngine, "assoc"}, "project1-assoc.golden"},
		{"project 1 - fields", []string{"--config=test/project1/corectl.yml ", connectToEngine, "fields"}, "project1-fields.golden"},
		{"project 1 - field numbers", []string{"--config=test/project1/corectl.yml ", connectToEngine, "field", "numbers"}, "project1-field-numbers.golden"},
		{"project 1 - eval", []string{"--config=test/project1/corectl.yml ", connectToEngine, "eval", "count(numbers)", "by", "xyz"}, "project1-eval-1.golden"},
		{"project 1 - objects", []string{"--config=test/project1/corectl.yml ", connectToEngine, "objects"}, "project1-objects.golden"},
		{"project 1 - data", []string{"--config=test/project1/corectl.yml ", connectToEngine, "--object", "my-hypercube", "data"}, "project1-data.golden"},
		{"project 1 - properties", []string{"--config=test/project1/corectl.yml ", connectToEngine, "--object", "my-hypercube", "properties"}, "project1-properties.golden"},
		{"project 1 - reload without progress", []string{"--config=test/project1/corectl.yml", connectToEngine, "reload", "--silent"}, "project1-reload-silent.golden"},

		// Project 2 has separate connections file
		{"project 2 - reload with connections", []string{connectToEngine, "-a=project2.qvf", "reload", "--script=test/project2/script.qvs", "--connections=test/project2/connections.yml", "--objects=test/project2/object-*.json"}, "project2-reload.golden"},
		{"project 2 - fields ", []string{"--config=test/project2/corectl.yml ", connectToEngine, "fields"}, "project2-fields.golden"},
		{"project 2 - data", []string{"--config=test/project2/corectl.yml ", connectToEngine, "--object", "my-hypercube-on-commandline", "data"}, "project2-data.golden"},

		{"project 3 - reload ", []string{"--config=test/project3/corectl.yml ", connectToEngine, "reload"}, "project3-reload.golden"},
		{"project 3 - fields", []string{"--config=test/project3/corectl.yml ", connectToEngine, "fields"}, "project3-fields.golden"},
		{"err 2", []string{connectToEngine, "--app=nosuchapp.qvf", "eval", "count(numbers)", "by", "xyz"}, "err-2.golden"},
		{"err 3", []string{connectToEngine, "--app=project1.qvf", "--object=nosuchobject", "data"}, "err-3.golden"},
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
		{"list apps", []string{connectToEngine, "apps"}, []string{"Name", "Last Reloaded", "ReadOnly", "Title", "project2.qvf", "project1.qvf"}},
		{"err 1", []string{"--engine=localhost:9999", "fields"}, []string{"Please check the --engine parameter or your config file", "Error details:  dial tcp"}},
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
