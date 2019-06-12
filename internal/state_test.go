package internal

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestBuildRestBaseUrl(t *testing.T) {
	assert.Equal(t, buildRestBaseURL("engine"), "http://engine")
	assert.Equal(t, buildRestBaseURL("engine.com"), "http://engine.com")
	assert.Equal(t, buildRestBaseURL("engine.com/app/appId"), "http://engine.com")
	assert.Equal(t, buildRestBaseURL("engine:1234"), "http://engine:1234")
	assert.Equal(t, buildRestBaseURL("ws://engine:1234"), "http://engine:1234")
	assert.Equal(t, buildRestBaseURL("wss://engine:1234"), "https://engine:1234")
	assert.Equal(t, buildRestBaseURL("wss://engine:1234/app/appId"), "https://engine:1234")
}

func TestBuildMetaUrl(t *testing.T) {
	assert.Equal(t, buildMetadataURL("engine", "appId"), "http://engine/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("engine:1234", "appId"), "http://engine:1234/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("wss://engine", "appId"), "https://engine/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("ws://engine", "appId"), "http://engine/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("wss://engine:1234", "appId"), "https://engine:1234/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("ws://engine:1234", "appId"), "http://engine:1234/v1/apps/appId/data/metadata")
	assert.Equal(t, buildMetadataURL("ws://engine:1234/app/appId", "appId"), "http://engine:1234/v1/apps/appId/data/metadata")
}

func TestBuildEngineUrl(t *testing.T) {
	assert.Equal(t, buildWebSocketURL("engine", "30"), "ws://engine/app/corectl/ttl/30")
	assert.Equal(t, buildWebSocketURL("engine:1234", "30"), "ws://engine:1234/app/corectl/ttl/30")
	assert.Equal(t, buildWebSocketURL("wss://engine", "30"), "wss://engine/app/corectl/ttl/30")
	assert.Equal(t, buildWebSocketURL("ws://engine", "30"), "ws://engine/app/corectl/ttl/30")
	assert.Equal(t, buildWebSocketURL("wss://engine:1234", "30"), "wss://engine:1234/app/corectl/ttl/30")
	assert.Equal(t, buildWebSocketURL("ws://engine:1234", "30"), "ws://engine:1234/app/corectl/ttl/30")
	assert.Equal(t, buildWebSocketURL("ws://engine:1234/sense/app/test.qvf", "30"), "ws://engine:1234/sense/app/test.qvf")
	assert.Equal(t, buildWebSocketURL("engine:1234/sense/app/test.qvf", "30"), "ws://engine:1234/sense/app/test.qvf")
	assert.Equal(t, buildWebSocketURL("ws://engine:1234/", "30"), "ws://engine:1234/")
}

func TestParseAppFromUrl(t *testing.T) {
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/app/test.qvf"), "test.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/app/test.qvf/"), "test.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/app/test.qvf/ttl/30"), "test.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/vp/app/test.qvf"), "test.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/vp/app/test.qvf/"), "test.qvf")
	assert.Equal(t, TryParseAppFromURL("engine/sense/vp/app/test.qvf/"), "test.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/vp/app/d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf/ttl/30"), "d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense/vp/app/d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf/ttl/30"), "d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf")
	assert.Equal(t, TryParseAppFromURL("ws://engine/"), "")
	assert.Equal(t, TryParseAppFromURL("ws://engine"), "")
	assert.Equal(t, TryParseAppFromURL("ws://engine/sense"), "")
}
