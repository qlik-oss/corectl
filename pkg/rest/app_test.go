package rest

import (
	"testing"
)

func TestIsAppID(t *testing.T) {
	cases := []struct {
		id  string
		exp bool
	}{
		{"", false},
		{"a-b-c-d-e", false},
		{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", true},
		{"6047f2b3-30b4-411e-baab-d3f60a79b95a", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"6047F2B3-30b4-411e-baab-d3f60a79b95a", false}, // Upper case.
		{"00000000-0000-0000-0000-00000000000", false},  // Wrong size.
		{"0000000-0000-0000-0000-000000000000", false},
		{"00000000-000-0000-0000-000000000000", false},
		{"00000000-0000-000-0000-000000000000", false},
		{"00000000-0000-0000-000-000000000000", false},
		{"blalba6047f2b3-30b4-411e-baab-d3f60a79b95ablabla", false},
	}
	for _, cas := range cases {
		if isAppID(cas.id) != cas.exp {
			t.Errorf("Expected %q to be %t", cas.id, cas.exp)
		}
	}
}

func TestIsItemID(t *testing.T) {
	cases := []struct {
		id  string
		exp bool
	}{
		{"", false},
		{"a-b", false},
		{"aaaaaaaaaaaaaaaaaaaaaaaa", true},
		{"000000000000000000000000", true},
		{"5ea0a15906d8b000019ba317", true},
		{"5eA0a15906d8b000019ba317", false}, // Upper case.
		{"aaaaaaaaaaa-aaaaaaaaaaaa", false}, // Contains hyphen.
		{"aaaaaaaaaaaaaaaaaaaaaaa", false},  // Wrong size.
		{"aaaaaaaaaaaaaaaaaaaaaaaaa", false},
	}
	for _, cas := range cases {
		if isItemID(cas.id) != cas.exp {
			t.Errorf("Expected %q to be %t", cas.id, cas.exp)
		}
	}
}
