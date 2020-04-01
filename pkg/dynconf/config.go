package dynconf

import (
	"fmt"
	"github.com/qlik-oss/corectl/pkg/log"
	leven "github.com/texttheater/golang-levenshtein/levenshtein"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var validProps = make(map[string]struct{})

func Glob(pattern string) []string {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		log.Warnf("Invalid glob pattern: %s", err)
	}
	if len(paths) == 0 {
		log.Warnf("No matches found for pattern %s", pattern)
	}
	return paths
}

//AddValidConfigFilePropertyName adds the supplied name to the list of valid properties that are allowed to appear in the main config file
func AddValidConfigFilePropertyName(name string) {
	validProps[name] = struct{}{}
}

// ReadYamlFile reads a yaml file from the supplied path and performs environment variable substitution.
func ReadYamlFile(path string, nameInLogs string, config interface{}) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("could not find %s '%s'\n", nameInLogs, path)
	}
	tempConfig := map[interface{}]interface{}{}
	err = yaml.Unmarshal(source, &tempConfig)
	if err != nil {
		log.Fatalf("invalid syntax in %s '%s': %s\n", nameInLogs, path, err)
	}
	err = subEnvVars(&tempConfig)
	if err != nil {
		log.Fatalf("bad substitution in '%s': %s\n", path, err)
	}
	if strConfig, err := convertMap(tempConfig); err == nil {
		reMarshal(strConfig, config)
	} else {
		log.Fatalf("could not parse %s '%s': %s\n", nameInLogs, path, err)
	}
}

// subEnvVars substitutes all the environment variables with their actual values in
// a map[string]interface{}, typically the unmarshallad yaml. (recursively)
func subEnvVars(m *map[interface{}]interface{}) error {
	for k, v := range *m {
		switch v.(type) {
		case string:
			envVar := v.(string)
			if strings.HasPrefix(envVar, "${") && strings.HasSuffix(envVar, "}") {
				envVar = strings.TrimSuffix(strings.TrimPrefix(envVar, "${"), "}")
				if val := os.Getenv(envVar); val != "" {
					(*m)[k] = val
				} else {
					return fmt.Errorf("environment variable '%s' not found", envVar)
				}
			}
		case map[interface{}]interface{}:
			m2 := v.(map[interface{}]interface{})
			if err := subEnvVars(&m2); err != nil {
				return err
			}
		}
	}
	return nil
}

// getSuggestion finds the best matching property within the specified Levenshtein distance limit
func getSuggestion(word string, validProps map[string]struct{}) string {
	op := leven.DefaultOptions // Default is cost 1 for del & ins, and 2 for substitution
	limit := 4
	min, suggestion := limit, ""
	for key := range validProps {
		dist := leven.DistanceForStrings([]rune(word), []rune(key), op)
		if dist < min {
			min = dist
			suggestion = key
		}
	}
	return suggestion
}

// convertMap turns {} -> {} map into string -> {} map
// returns error if non-string was present in input map
func convertMap(m map[interface{}]interface{}) (map[string]interface{}, error) {
	strMap := map[string]interface{}{}
	for k, v := range m {
		if s, ok := k.(string); ok {
			strMap[s] = v
		} else {
			return strMap, fmt.Errorf("property '%v' is not a string", k)
		}
	}
	return strMap, nil
}

func reMarshal(m map[string]interface{}, ref interface{}) error {
	bytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, ref)
	if err != nil {
		return err
	}
	return nil
}

// findConfigFile finds a file with the given fileName with yml or yaml extension.
// Returns absolute path
func findConfigFile(fileName string) string {
	configFile := ""
	if _, err := os.Stat(fileName + ".yml"); !os.IsNotExist(err) {
		configFile = fileName + ".yml"
	} else if _, err := os.Stat(fileName + ".yaml"); !os.IsNotExist(err) {
		configFile = fileName + ".yaml"
	}
	if configFile != "" {
		absConfig, err := filepath.Abs(configFile) // Convert to abs path
		if err != nil {
			log.Fatalf("unexpected error when converting to absolute filepath: %s\n", err)
		}
		configFile = absConfig
	}
	return configFile
}

// validateProps checks if there are unknown properties in the config
// configPath is passed for error logging purposes.
func validateProps(config map[interface{}]interface{}, configPath string) {
	invalidProps := []string{}
	suggestions := map[string]string{}
	for key := range config {
		keyString, ok := key.(string)
		if !ok {
			// If there is a non-string in the yaml, this will surely be an invalid props
			keyString = fmt.Sprint(key)
		}
		if _, ok := validProps[keyString]; !ok {
			if suggestion := getSuggestion(keyString, validProps); suggestion != "" {
				suggestions[keyString] = suggestion
			} else {
				invalidProps = append(invalidProps, fmt.Sprintf("'%s'", keyString)) // For pretty printing
			}
		}
	}
	if len(invalidProps)+len(suggestions) > 0 {
		errorMessage := []string{}
		errorMessage = append(errorMessage,
			fmt.Sprintf("corectl found invalid properties when validating the config file '%s'.", configPath))
		for key, value := range suggestions {
			errorMessage = append(errorMessage, fmt.Sprintf("  '%s': did you mean '%s'?", key, value))
		}
		if len(invalidProps) > 0 {
			prepend := "M" // Capitalize M if there were no suggestions
			if len(suggestions) > 0 {
				prepend = "Also, m" // Add also if there were suggestions
			}
			errorMessage = append(errorMessage,
				fmt.Sprintf("%systerious properties: %s", prepend, strings.Join(invalidProps, ", ")))
		}
		log.Fatalln(strings.Join(errorMessage, "\n"))
	}
}

func mergeContext(config *map[interface{}]interface{}, contextName string) {
	contextHandler := NewContextHandler()

	if contextName == "" {
		contextName = contextHandler.Current
	}

	context := contextHandler.Get(contextName)

	if context == nil {
		return
	}

	log.Verboseln("Merging config with context: " + contextName)

	for k, v := range context {
		if _, ok := (*config)[k]; ok {
			if k == "headers" { // Handle headers separately.
				continue
			} else {
				log.Warnf("Property '%s' exists in both current context and config, using property from config\n", k)
			}
		} else {
			(*config)[k] = v
		}
	}

	ctxHeaders := context.Headers()
	if value, ok := (*config)["headers"]; ok {
		cfgHeaders := convertToStringStringMap(value)
		mergeHeaders(cfgHeaders, ctxHeaders)
		(*config)["headers"] = cfgHeaders
	} else {
		(*config)["headers"] = ctxHeaders
	}
}
