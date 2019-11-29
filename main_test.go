package main

import "testing"

import "bytes"

func TestAbs(t *testing.T) {
	input := bytes.NewBufferString("++.++")
	got, _ := eval(input)
	want := "4"

	if got != want {
		t.Errorf("want: '%s', got: '%s'", want, got)
	}
}
