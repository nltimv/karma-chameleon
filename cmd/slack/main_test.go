package main

import "testing"

func TestPass(t *testing.T) {
	expected := 5
	actual := 2 + 3

	if expected != actual {
		t.Errorf("got %q, wanted %q", actual, expected)
	}
}
