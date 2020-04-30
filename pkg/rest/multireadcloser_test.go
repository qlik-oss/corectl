package rest_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/qlik-oss/corectl/pkg/rest"
)

type readCloser struct {
	io.Reader
	closed bool
}

func NewReadCloser(r io.Reader) *readCloser {
	return &readCloser{r, false}
}

func (rc *readCloser) Close() error {
	if rc.closed {
		return fmt.Errorf("already closed")
	}
	rc.closed = true
	return nil
}

func TestMultiReadCloser(t *testing.T) {
	strs := []string{"Hello", " dear", " reader."}
	readClosers := make([]io.ReadCloser, len(strs))
	for i, s := range strs {
		buf := bytes.NewBuffer([]byte(s))
		readClosers[i] = NewReadCloser(buf)
	}

	mrc := rest.MultiReadCloser(readClosers...)
	b, err := ioutil.ReadAll(mrc)
	if err != nil {
		t.Fatal(err)
	}
	if s := strings.Join(strs, ""); string(b) != s {
		t.Errorf("Expected %q but got %q", s, string(b))
	}
	mrc.Close()
	for _, rc := range readClosers {
		if !rc.(*readCloser).closed {
			t.Error("ReadCloser was not closed!")
		}
	}
}
