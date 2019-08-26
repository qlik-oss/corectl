package internal

import (
	"encoding/json"
)

const baseSchemaURL string = "https://deploy-preview-514--core-website.netlify.com/json-schema/"

// InjectSchemaIntoProperties adds a '$schema' property to the object with the url to the json definition for IntelliSense usage
func InjectSchemaIntoProperties(properties json.RawMessage, entityType string) json.RawMessage {
	var customProperties map[string]interface{}
	json.Unmarshal(properties, &customProperties)

	// only inject a schema property if one does not already exist
	if customProperties["$schema"] == nil {
		customProperties["$schema"] = baseSchemaURL + "generic-" + entityType + ".json"
		res, _ := json.Marshal(customProperties)
		return res
	}
	return properties
}
