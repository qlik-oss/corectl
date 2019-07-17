package internal

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var source = []byte(`test: ${_TEST1_}
nested:
  nest2:
    name: foo
    name2: foo
    test2: ${_TEST2_}
  test3: ${_TEST3_}
list:
  - ${_TEST4_}
`)

func TestEnvVarSubstitution(t *testing.T) {
	config := readConfig(source)
	fmt.Println(config)
	err := subEnvVars(config)
	assert.Error(t, err)
	os.Setenv("_TEST1_", "TEST1")
	err = subEnvVars(readConfig(source))
	assert.Error(t, err)
	os.Setenv("_TEST2_", "TEST2")
	err = subEnvVars(readConfig(source))
	assert.Error(t, err)
	os.Setenv("_TEST3_", "TEST3")
	err = subEnvVars(readConfig(source))
	// We don't substitute env-variables in lists
	assert.Nil(t, err)
}

func readConfig(source []byte) (config *map[interface{}]interface{}) {
	config = &(map[interface{}]interface{}{})
	if err := yaml.Unmarshal(source, config); err != nil {
		fmt.Println(err)
		return nil
	}
	return
}
