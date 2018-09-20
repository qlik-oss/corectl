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

func TestCorectl(t *testing.T) {
	connectToEngine := "--engine=" + *engineIP
	tests := []struct {
		name   string
		args   []string
		golden string
	}{
		{"reload project1", []string{"--config=test/project1/qli.yml", connectToEngine, "reload"}, "project1-reload.golden"},
		{"tables project 1", []string{"--config=test/project1/qli.yml ", connectToEngine, "tables"}, "project1-tables.golden"},
		{"assoc project 1", []string{"--config=test/project1/qli.yml ", connectToEngine, "assoc"}, "project1-assoc.golden"},
		{"fields project 1", []string{"--config=test/project1/qli.yml ", connectToEngine, "fields"}, "project1-fields.golden"},
		{"field numbers project 1", []string{"--config=test/project1/qli.yml ", connectToEngine, "field", "numbers"}, "project1-field-numbers.golden"},
		{"eval project 1", []string{"--config=test/project1/qli.yml ", connectToEngine, "eval", "count(numbers)", "by", "xyz"}, "project1-eval-1.golden"},

		{"reload project 2", []string{"--config=test/project2/qli.yml ", connectToEngine, "reload"}, "project2-reload.golden"},
		{"fields project 2", []string{"--config=test/project2/qli.yml ", connectToEngine, "fields"}, "project2-fields.golden"},
		{"reload project 2 with connections", []string{connectToEngine, "-a project2.qvf ", "reload", "--script=test/project2/script.qvs", "--connections=test/project2/connections.yml"}, "project2-reload.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
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
