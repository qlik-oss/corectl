package internal

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

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
