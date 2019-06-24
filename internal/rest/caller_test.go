package rest

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestBuildRestBaseUrl(t *testing.T) {
	f := func(s string) string {
		return buildRestBaseURL(s).String()
	}

	assert.Equal(t, f("engine"), "http://engine")
	assert.Equal(t, f("engine.com"), "http://engine.com")
	assert.Equal(t, f("engine.com/app/appId"), "http://engine.com")
	assert.Equal(t, f("engine:1234"), "http://engine:1234")
	assert.Equal(t, f("ws://engine:1234"), "http://engine:1234")
	assert.Equal(t, f("wss://engine:1234"), "https://engine:1234")
	assert.Equal(t, f("wss://engine:1234/app/appId"), "https://engine:1234")
}
