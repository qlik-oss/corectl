package dynconf_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/qlik-oss/corectl/pkg/boot"
	"github.com/qlik-oss/corectl/pkg/dynconf"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var cmd *cobra.Command
var configFile string

var context = map[string]interface{}{
	"verbose": false,
	"traffic": true,
	"headers": map[string]string{
		"accept":        "text/plain",
		"Authorization": "bla",
	},
}

var config = map[string]interface{}{
	"verbose": true,
	"no-data": false,
	"headers": map[string]string{
		"Accept": "encoding/json",
		"Cookie": "MONSTER",
		"date":   "tomorrow",
	},
}

var commandLine = map[string]interface{}{
	"verbose": true,
	"headers": map[string]string{
		"Accept": "image/webp",
		"cookie": "gingerbread",
	},
}

func TestMain(m *testing.M) {
	dir := setupTempDir()
	cmd = &cobra.Command{}
	setupContext(dir)
	setupConfig(dir)
	boot.InjectGlobalFlags(cmd, false)
	boot.InjectAppWebSocketFlags(cmd, false)
	cmd.Execute() // Execute to parse the flag settings.
	code := m.Run()
	err := os.RemoveAll(dir)
	if err != nil {
		fmt.Println("failed to remove temp dir:", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func setupContext(dir string) {
	dynconf.ContextFilePath = path.Join(dir, "contexts.yml")
	dynconf.CreateContext("test", context)
}

func setupConfig(dir string) {
	configFile = path.Join(dir, "config.yml")
	b, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(configFile, b, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestDefault(t *testing.T) {
	// Nothing is set, should all be defaults.
	settings := dynconf.ReadSettings(cmd)
	if settings.GetBool("verbose") {
		t.Error("expected verbose to be false")
	}
	if !settings.IsUsingDefaultValue("verbose") {
		t.Error("verbose to have default value")
	}
}

func TestContext(t *testing.T) {
	// Nothing is set, should all be defaults.
	cmd.Flags().Set("context", "test")
	settings := dynconf.ReadSettings(cmd)
	if key := "verbose"; settings.GetBool(key) != context[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if key := "traffic"; settings.GetBool(key) != context[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if settings.IsUsingDefaultValue("verbose") {
		t.Error("verbose to not be default value, it is set in the context")
	}
}

func TestConfig(t *testing.T) {
	cmd.Flags().Set("context", "test")
	cmd.Flags().Set("config", configFile)
	settings := dynconf.ReadSettings(cmd)
	if key := "verbose"; settings.GetBool(key) != config[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if key := "traffic"; settings.GetBool(key) != context[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if key := "no-data"; settings.GetBool(key) != config[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if settings.IsUsingDefaultValue("verbose") {
		t.Error("verbose to not be default value, it is set in the context")
	}
	headers := settings.GetHeaders()
	cfgHeaders := config["headers"].(map[string]string)
	for k, v := range cfgHeaders {
		if headers.Get(k) != v {
			t.Errorf("was %q expected %q", headers.Get(k), v)
		}
	}
}

func TestCommandLine(t *testing.T) {
	for k, v := range commandLine {
		switch value := v.(type) {
		case map[string]string:
			header := value
			for hk, hv := range header {
				cmd.Flags().Set(k, hk+"="+hv)
			}
		default:
			cmd.Flags().Set(k, fmt.Sprint(value))
		}
	}
	cmd.Flags().Set("context", "test")
	cmd.Flags().Set("config", configFile)
	settings := dynconf.ReadSettings(cmd)
	if key := "verbose"; settings.GetBool(key) != commandLine[key] {
		t.Errorf("expected %q to %t", key, config[key])
	}
	if key := "traffic"; settings.GetBool(key) != context[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if key := "no-data"; settings.GetBool(key) != config[key] {
		t.Errorf("expected %q to %t", key, context[key])
	}
	if settings.IsUsingDefaultValue("verbose") {
		t.Error("verbose to not be default value, it is set in the context")
	}
	headers := settings.GetHeaders()
	// Copy command-line headers into expected
	expected := map[string]string{}
	for k, v := range commandLine["headers"].(map[string]string) {
		expected[k] = v
	}
	// Set expected key-value pairs from context and config
	expected["Authorization"] = context["headers"].(map[string]string)["Authorization"]
	expected["date"] = config["headers"].(map[string]string)["date"]
	for k, v := range expected {
		if headers.Get(k) != v {
			t.Errorf("header %q was %q, expected %q", k, headers.Get(k), v)
		}
	}
	headerLength := len(map[string][]string(headers))
	if headerLength != len(expected) {
		t.Errorf("header length was %d, expected %d", headerLength, len(expected))
	}
}

func setupTempDir() string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = os.Mkdir(path.Join(dir, dynconf.ContextDir), 0755)
	return dir
}
