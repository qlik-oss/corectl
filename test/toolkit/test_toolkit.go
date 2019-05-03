package toolkit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var createdApps []string

type Params struct {
	T       *testing.T
	App     string
	Config  string
	Engine  string
	Headers string
	NoData  string
	Traffic string
	Ttl     string
	Verbose string

	expectGolden   bool
	expectIncludes []string
	expectError    bool
	expectOK       bool
	expectEqual    string
}

func (p *Params) ExpectGolden() *Params {
	pc := *p // Shallow clone
	pc.expectGolden = true
	return &pc
}
func (p *Params) ExpectOK() *Params {
	pc := *p // Shallow clone
	pc.expectOK = true
	return &pc
}
func (p *Params) ExpectError() *Params {
	pc := *p // Shallow clone
	pc.expectError = true
	return &pc
}
func (p *Params) ExpectIncludes(items ...string) *Params {
	pc := *p // Shallow clone
	pc.expectIncludes = items
	return &pc
}

func (p *Params) ExpectEqual(item string) *Params {
	pc := *p // Shallow clone
	pc.expectEqual = item
	return &pc
}

func (p *Params) Describe(title string) *Params {
	p.T.Log(title)
	return p
}

func toGoldenFileName(name string) string {
	goldenBaseName := strings.Replace(name, "/", "_", -1)
	return goldenBaseName + ".golden"
}

func (p *Params) WithParams(newP Params) Params {
	pc := *p // Shallow clone
	if newP.Config != "" {
		pc.Config = newP.Config
	}
	if newP.Engine != "" {
		pc.Engine = newP.Engine
	}
	if newP.App != "" {
		pc.App = newP.App
	}
	if newP.Verbose != "" {
		pc.Verbose = newP.Verbose
	}
	if newP.Traffic != "" {
		pc.Traffic = newP.Traffic
	}
	if newP.Ttl != "" {
		pc.Ttl = newP.Ttl
	}
	if newP.Headers != "" {
		pc.Headers = newP.Headers
	}
	if newP.NoData != "" {
		pc.NoData = newP.NoData
	}
	return pc
}

func (p *Params) Run(command ...string) []byte {
	var output []byte
	name := strings.Join(command, " ")
	p.T.Run(name, func(t *testing.T) {
		var args []string
		if p.App != "" {
			args = append(args, "--app", p.App)
		}
		if p.Engine != "" {
			args = append(args, "--engine", p.Engine)
		}
		if p.Config != "" {
			args = append(args, "--config", p.Config)
		}
		if p.Verbose == "true" {
			args = append(args, "--verbose")
		}
		if p.Headers != "" {
			args = append(args, "--headers", "\""+p.Headers+"\"")
		}
		if p.NoData == "true" {
			args = append(args, "--no-data")
		}
		if p.Traffic == "true" {
			args = append(args, "--traffic")
		}
		if p.Ttl != "" {
			args = append(args, "--ttl", p.Ttl)
		}
		args = append(args, command...)

		createdApp := createsApp(args)

		if createdApp != "" {
			createdApps = append(createdApps, createdApp)
		}

		cmd := exec.Command(binaryPath, args...)

		goldenName := toGoldenFileName(t.Name())

		t.Log("\u001b[35m Executing command:" + strings.Join(cmd.Args, " ") + "\u001b[0m")
		var err error
		output, err = cmd.CombinedOutput()

		actual := string(output)
		t.Log("\u001b[35m Output:\n" + actual)
		if p.expectOK {
			if err != nil {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, false, err != nil, err)
			}
		} else if p.expectError {
			if err == nil {
				t.Fatalf("%s\nexpected (err == nil) to be %v, but got %v. err: %v", output, false, err == nil, err)
			}
		} else if p.expectGolden {
			golden := newGoldenFile(t, goldenName)

			if update {
				golden.write(actual)
			}
			expected := golden.load()

			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("diff: %v", diff(expected, actual))
			}
		} else if p.expectEqual != "" {
			if err != nil {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, false, err != nil, err)
			}
			if strings.Trim(actual, " \t\n") != strings.Trim(p.expectEqual, " \t\n") {
				t.Fatalf("Output did not equal string '%v'\nReceived:\n%v", p.expectEqual, actual)
			}
		} else if p.expectIncludes != nil && len(p.expectIncludes) > 0 {
			if err != nil {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, false, err != nil, err)
			}
			for _, sub := range p.expectIncludes {
				if !strings.Contains(actual, sub) {
					t.Fatalf("Output did not contain substring '%v'\nReceived:\n%v", sub, actual)
				}
			}
		}
	})
	return output
}

func createsApp(args []string) string {
	var buildFound bool
	var nextIsApp bool
	var app string
	for _, item := range args {
		if item == "build" {
			buildFound = true
		} else if nextIsApp {
			app = item
			nextIsApp = false
		} else if item == "--app" || item == "-a" {
			nextIsApp = true
		} else if strings.Index(item, "--app=") == 0 || strings.Index(item, "-a=") == 0 {
			app = strings.Split(item, "=")[1]
		}
	}
	if buildFound {
		return app
	}
	return ""
}

func (p *Params) Reset() {
	if createdApps != nil {
		for _, appId := range createdApps {
			p.ExpectOK().Run("app", "rm", appId, "--suppress")
		}
		createdApps = []string{}

	} else {
		p.T.Log("No apps found when resetting")
	}
}

func buildCorectl() {
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
}

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "corectl.exe"
	}

	return "corectl"
}

var binaryName = getBinaryName()

var binaryPath string

func GetTestFilePath() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}

func init() {
	fmt.Println("RUNNING TEST MAIN")
	buildCorectl()
	os.Setenv("CORECTL_TEST_CONNECT", "corectl-test-connector")
	os.Setenv("ENGINE_URL", "localhost:9076")

}
