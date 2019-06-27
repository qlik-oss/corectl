package internal

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/spf13/viper"
)

func TestParseEngineURL1(t *testing.T) {
	testURL(t, "engine", map[string]string{
		"Scheme": "ws",
		"Host": "engine",
		"Path": "",
	}, true)
	testURL(t, "engine:1234", map[string]string{
		"Scheme": "ws",
		"Host": "engine:1234",
		"Path": "",
	}, true)
	testURL(t, "localhost:1234", map[string]string{
		"Scheme": "ws",
		"Host": "localhost:1234",
		"Path": "",
	}, true)
	testURL(t, "localhost:1234/app/apa", map[string]string{
		"Scheme": "ws",
		"Host": "localhost:1234",
		"Path": "/app/apa",
	}, true)
	testURL(t, "engine/path", map[string]string{
		"Scheme": "ws",
		"Host": "engine",
		"Path": "/path",
	}, true)
	testURL(t, "engine.com", map[string]string{
		"Scheme": "ws",
		"Host": "engine.com",
		"Path": "",
	}, true)
	testURL(t, "http://engine.com", map[string]string{
		"Scheme": "ws",
		"Host": "engine.com",
		"Path": "",
	}, true)
	testURL(t, "http://engine/path", map[string]string{
		"Scheme": "ws",
		"Host": "engine",
		"Path": "/path",
	}, true)
	testURL(t, "http://127.0.0.1:1234", map[string]string{
		"Scheme": "ws",
		"Host": "127.0.0.1:1234",
		"Path": "",
	}, true)
	testURL(t, "127.0.0.1:1234", map[string]string{
		"Scheme": "ws",
		"Host": "127.0.0.1:1234",
		"Path": "",
	}, true)
	testURL(t, "127.0.0.1", map[string]string{
		"Scheme": "ws",
		"Host": "127.0.0.1",
		"Path": "",
	}, true)
	testURL(t, "127.0.0.1/app/foo", map[string]string{
		"Scheme": "ws",
		"Host": "127.0.0.1",
		"Path": "/app/foo",
	}, true)
	testURL(t, "ws://localhost:1234/app/foo", map[string]string{
		"Scheme": "ws",
		"Host": "localhost:1234",
		"Path": "/app/foo",
	}, true)
}

func testURL(t *testing.T, s string, fields map[string]string, pass bool) (u *url.URL) {
	u, err := parseEngineURL(s)
	if pass {
		if !assert.Nil(t, err) {
			return
		}
	} else {
		assert.Error(t, err)
		return
	}
	v := reflect.ValueOf(*u)
	for f, expected := range(fields) {
		fval := string(v.FieldByName(f).String())
		s_exp := fmt.Sprintf("'%s'= %s", f, expected)
		s_fval := fmt.Sprintf("'%s'= %s", f, fval)
		assert.Equal(t, s_exp, s_fval)
	}
	return
}

func TestParseEngineUrl2(t *testing.T) {
	// Wrapper function
	f := func(s string) string {
		viper.Set("engine", s)
		return ParseEngineURL().String()
	}
	assert.Equal(t, "ws://engine", f("engine"))
	assert.Equal(t, "ws://engine:1234", f("engine:1234"))
	assert.Equal(t, "ws://engine:1234/app/foo", f("engine:1234/app/foo"))
	assert.Equal(t, "ws://engine:1234/app/foo.qvf", f("engine:1234/app/foo.qvf"))
	assert.Equal(t, "ws://127.0.0.1:1234/app/foo.qvf", f("127.0.0.1:1234/app/foo.qvf"))
	assert.Equal(t, "ws://engine", f("http://engine"))
}

func TestBuildEngineUrl(t *testing.T) {
	// Wrapper function
	f := func(s, ttl string) string {
		viper.Set("engine", s)
		return buildWebSocketURL(ttl)
	}
	assert.Equal(t, "ws://engine/app/corectl/ttl/30", f("engine", "30"))
	assert.Equal(t, "ws://engine:1234/app/corectl/ttl/30", f("engine:1234", "30"))
	assert.Equal(t, "wss://engine/app/corectl/ttl/30", f("wss://engine", "30"))
	assert.Equal(t, "ws://engine/app/corectl/ttl/30", f("ws://engine", "30"))
	assert.Equal(t, "wss://engine:1234/app/corectl/ttl/30", f("wss://engine:1234", "30"))
	assert.Equal(t, "ws://engine:1234/app/corectl/ttl/30", f("ws://engine:1234", "30"))
	assert.Equal(t, "ws://engine:1234/sense/app/test.qvf", f("ws://engine:1234/sense/app/test.qvf", "30"))
	assert.Equal(t, "ws://engine:1234/sense/app/test.qvf", f("engine:1234/sense/app/test.qvf", "30"))
	assert.Equal(t, "ws://engine:1234/", f("ws://engine:1234/", "30"))
	assert.Equal(t, "ws://engine:1234/", f("http://engine:1234/", "30"))
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
