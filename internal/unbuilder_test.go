package internal

import (
	"fmt"
	"testing"
)

func TestBuildName(t *testing.T) {
	fmt.Println(buildEntityFilename("wefwef", "mastesrobject", "table", "'='Halleluljah moment'"))
}

func TestUnbuildRegex(t *testing.T) {
	// Inline function to avoid scope pollution
	find := func(s string) []string {
		return matchAllNonAlphaNumeric.FindAllString(s, -1)
	}
	pass := []string{
		"объект", "αντικείμενο", "חפץ",
		"対象", "object", "åbjäkt", "æble",
		"___sys", "αντικείμενο19", "hell-here",
		"1982", "аяАЯЁё",
	}
	fail := []string{
		"@object", "$erver", "\"hello\"",
		"p:;", "~som", "|other", "=thing",
		"hello there",
	}
	for _, s := range pass {
		if found := find(s); len(found) != 0 {
			t.Errorf("%s should not match the regex but found %v", s, found)
		}
	}
	for _, s := range fail {
		if found := find(s); len(found) == 0 {
			t.Errorf("%s should match the regex", s)
		}
	}
}
