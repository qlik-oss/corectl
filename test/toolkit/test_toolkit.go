package toolkit

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var (
	createdApps   []string
	goldPolishers []func(string) string
)

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

	expectGolden               bool
	expectIncludes             []string
	exectJsonArrayPropertyName string
	expectJsonArrayValues      []string
	expectError                bool
	expectOK                   bool
	expectEqual                string
}

func AddGoldPolisher(from string, to string) {
	polisher := func(content string) string {
		newConnectionCreatedRegexp := regexp.MustCompile("(?m)" + from)
		return newConnectionCreatedRegexp.ReplaceAllString(content, to)
	}
	goldPolishers = append(goldPolishers, polisher)
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

func (p *Params) ExpectEmptyJsonArray() *Params {
	pc := *p // Shallow clone
	pc.exectJsonArrayPropertyName = "NA"
	pc.expectJsonArrayValues = []string{}
	return &pc
}
func (p *Params) ExpectJsonArray(propertyName string, items ...string) *Params {
	pc := *p // Shallow clone
	pc.exectJsonArrayPropertyName = propertyName
	pc.expectJsonArrayValues = items
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

func (p *Params) filterForGold(content string) string {
	for _, x := range p.goldPolishers {
		content = x(content)
	}
	return content
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
		}

		if p.expectGolden {
			golden := newGoldenFile(t, goldenName)

			actualFiltered := p.filterForGold(actual)

			if update {
				golden.write(actualFiltered)
			}
			expected := golden.load()

			if !reflect.DeepEqual(expected, actualFiltered) {
				t.Fatalf("diff: %v", diff(expected, actualFiltered))
			}
		} else if p.expectEqual != "" {
			if err != nil {
				t.Fatalf("%s\nexpected (err != nil) to be %v, but got %v. err: %v", output, false, err != nil, err)
			}
			if strings.Trim(actual, " \t\n") != strings.Trim(p.expectEqual, " \t\n") {
				t.Fatalf("Output did not equal string '%v'\nReceived:\n%v", p.expectEqual, actual)
			}
		} else if p.expectJsonArrayValues != nil {
			var jsonArray []map[string]string
			json.Unmarshal(output, &jsonArray)
			fmt.Println(len(jsonArray))
			fmt.Println(len(p.expectJsonArrayValues))
			assert.Equal(t, len(jsonArray), len(p.expectJsonArrayValues), "Wrong size of array")
			for i, value := range p.expectJsonArrayValues {
				assert.Equal(t, value, jsonArray[i][p.exectJsonArrayPropertyName], "Unexpected value in json array")
			}
		} else if p.expectIncludes != nil && len(p.expectIncludes) > 0 {
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
			p.Run("app", "rm", appId, "--suppress")
		}
		createdApps = []string{}

	} else {
		p.T.Log("No apps found when resetting")
	}
}
