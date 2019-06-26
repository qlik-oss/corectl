package rest

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test is more like documentation of the golang's net/url
func TestURL(t *testing.T) {
	testit(t, "engine", map[string]string{
		"Scheme": "",
		"Host": "",
		"Path": "engine",
	}, true)
	testit(t, "engine:1234", map[string]string{
		"Scheme": "engine",
		"Host": "",
		"Opaque": "1234",
		"Path": "",
	}, true)
	testit(t, "localhost:1234", map[string]string{
		"Scheme": "localhost",
		"Host": "",
		"Opaque": "1234",
		"Path": "",
	}, true)
	testit(t, "localhost:1234/app/apa", map[string]string{
		"Scheme": "localhost",
		"Host": "",
		"Opaque": "1234/app/apa",
		"Path": "",
	}, true)
	testit(t, "engine/path", map[string]string{
		"Scheme": "",
		"Host": "",
		"Path": "engine/path",
	}, true)
	testit(t, "engine.com", map[string]string{
		"Scheme": "",
		"Host": "",
		"Path": "engine.com",
	}, true)
	testit(t, "http://engine.com", map[string]string{
		"Scheme": "http",
		"Host": "engine.com",
		"Path": "",
	}, true)
	testit(t, "http://engine/path", map[string]string{
		"Scheme": "http",
		"Host": "engine",
		"Path": "/path",
	}, true)
	testit(t, "http://127.0.0.1:1234", map[string]string{
		"Scheme": "http",
		"Host": "127.0.0.1:1234",
		"Path": "",
	}, true)
	testit(t, "127.0.0.1:1234", map[string]string{}, false)
	testit(t, "127.0.0.1", map[string]string{
		"Scheme": "",
		"Host": "",
		"Path": "127.0.0.1",
	}, true)
	testit(t, "127.0.0.1/app/foo", map[string]string{
		"Scheme": "",
		"Host": "",
		"Path": "127.0.0.1/app/foo",
	}, true)
	u := testit(t, "ws://localhost:1234/app/foo", map[string]string{
		"Scheme": "ws",
		"Host": "localhost:1234",
		"Path": "/app/foo",
	}, true)
	u.Path = ""
	assert.Equal(t, "ws://localhost:1234", u.String())
}

func testit(t *testing.T, s string, fields map[string]string, pass bool) (u *url.URL) {
	u, err := url.Parse(s)
	if pass {
		if !assert.Nil(t, err) {
			return
		}
	} else {
		assert.Error(t, err)
		return
	}
	assert.Equal(t, u.String(), s)
	v := reflect.ValueOf(*u)
	for f, expected := range(fields) {
		fval := string(v.FieldByName(f).String())
		s_exp := fmt.Sprintf("'%s'= %s", f, expected)
		s_fval := fmt.Sprintf("'%s'= %s", f, fval)
		assert.Equal(t, s_exp, s_fval)
	}
	return
}
