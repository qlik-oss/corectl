package internal

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestBuildMetaUrl(t *testing.T) {
	assert.Equal(t, buildMetadataURL("engine", "appId"), "http://engine:9076/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("engine:1234", "appId"), "http://engine:1234/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("wss://engine", "appId"), "https://engine:9076/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("ws://engine", "appId"), "http://engine:9076/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("wss://engine:1234", "appId"), "https://engine:1234/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("ws://engine:1234", "appId"), "http://engine:1234/v1/apps/appId/data/metadata")
}

func TestBuildEngineUrl(t *testing.T) {
	assert.Equal(t, buildWebSocketURL("engine", "30"), "ws://engine:9076/app/engineData/ttl/30")
	assert.Equal(t, buildWebSocketURL("engine:1234", "30"), "ws://engine:1234/app/engineData/ttl/30")
	assert.Equal(t, buildWebSocketURL("wss://engine", "30"), "wss://engine:9076/app/engineData/ttl/30")
	assert.Equal(t, buildWebSocketURL("ws://engine", "30"), "ws://engine:9076/app/engineData/ttl/30")
	assert.Equal(t, buildWebSocketURL("wss://engine:1234", "30"), "wss://engine:1234/app/engineData/ttl/30")
	assert.Equal(t, buildWebSocketURL("ws://engine:1234", "30"), "ws://engine:1234/app/engineData/ttl/30")
}
