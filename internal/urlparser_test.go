package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEngineUrl(t *testing.T) {
	// Wrapper function for calling multireturn function
	f := func(s string) string {
		return ParseEngineURL(s).String()
	}
	assert.Equal(t, "ws://engine", f("engine"))
	assert.Equal(t, "ws://engine:1234", f("engine:1234"))
	assert.Equal(t, "ws://engine:1234/app/foo", f("engine:1234/app/foo"))
	assert.Equal(t, "ws://engine:1234/app/foo.qvf", f("engine:1234/app/foo.qvf"))
	assert.Equal(t, "ws://127.0.0.1:1234/app/foo.qvf", f("127.0.0.1:1234/app/foo.qvf"))
	assert.Equal(t, "ws://engine", f("http://engine"))
}

func TestBuildEngineUrl(t *testing.T) {
	assert.Equal(t, "ws://engine/app/corectl/ttl/30", buildWebSocketURL("engine", "30"))
	assert.Equal(t, "ws://engine:1234/app/corectl/ttl/30", buildWebSocketURL("engine:1234", "30"))
	assert.Equal(t, "wss://engine/app/corectl/ttl/30", buildWebSocketURL("wss://engine", "30"))
	assert.Equal(t, "ws://engine/app/corectl/ttl/30", buildWebSocketURL("ws://engine", "30"))
	assert.Equal(t, "wss://engine:1234/app/corectl/ttl/30", buildWebSocketURL("wss://engine:1234", "30"))
	assert.Equal(t, "ws://engine:1234/app/corectl/ttl/30", buildWebSocketURL("ws://engine:1234", "30"))
	assert.Equal(t, "ws://engine:1234/sense/app/test.qvf", buildWebSocketURL("ws://engine:1234/sense/app/test.qvf", "30"))
	assert.Equal(t, "ws://engine:1234/sense/app/test.qvf", buildWebSocketURL("engine:1234/sense/app/test.qvf", "30"))
	assert.Equal(t, "ws://engine:1234/", buildWebSocketURL("ws://engine:1234/", "30"))
	assert.Equal(t, "ws://engine:1234/", buildWebSocketURL("http://engine:1234/", "30"))
}

func TestParseAppFromUrl(t *testing.T) {
	assert.Equal(t, "test.qvf", TryParseAppFromURL("ws://engine/sense/app/test.qvf"))
	assert.Equal(t, "test.qvf", TryParseAppFromURL("ws://engine/sense/app/test.qvf/"))
	assert.Equal(t, "test.qvf", TryParseAppFromURL("ws://engine/sense/app/test.qvf/ttl/30"))
	assert.Equal(t, "test.qvf", TryParseAppFromURL("ws://engine/sense/vp/app/test.qvf"))
	assert.Equal(t, "test.qvf", TryParseAppFromURL("ws://engine/sense/vp/app/test.qvf/"))
	assert.Equal(t, "test.qvf", TryParseAppFromURL("engine/sense/vp/app/test.qvf/"))
	assert.Equal(t, "d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf", TryParseAppFromURL("ws://engine/sense/vp/app/d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf/ttl/30"))
	assert.Equal(t, "d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf", TryParseAppFromURL("ws://engine/sense/vp/app/d6c281c1-3463-4b0a-8251-ed747e9e426e.qvf/ttl/30"))
	assert.Equal(t, "", TryParseAppFromURL("ws://engine/"))
	assert.Equal(t, "", TryParseAppFromURL("ws://engine"))
	assert.Equal(t, "", TryParseAppFromURL("ws://engine/sense"))
}
