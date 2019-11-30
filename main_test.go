package main

import (
	"bytes"
	"os"
	"testing"
)

func eval(s string) string {
	bf := newBf(bytes.NewBufferString(s))

	bf.Run()

	return bf.IntString()
}

func TestIntData(t *testing.T) {
	type TestCase struct {
		Input string
		Want  string
	}

	for _, tc := range []TestCase{
		{"++.+>[-.]", "30"},
		{"+", "1"},
		{"[ empty [loop] inner ] >> ++", "002"},
		{"[ [] ++] >+", "01"},
		{"+++[>>+<<-  ]+>++", "123"},
		{"++++++++ >++++ >++   >+++   >+++     >+ >+++", "8423313"},
		{` ++ > +++++  [ < + > - ] ++++ ++++ [ < +++ +++ > -] < .`, "550"},
	} {
		got := eval(tc.Input)

		if tc.Want != got {
			t.Errorf("Fail on input: %s\nWant: '%s'\nGot : '%s'", tc.Input, tc.Want, got)
		}
	}
}

func TestHelloWorld(t *testing.T) {
	f, _ := os.Open("hello_world.bf")

	bf := newBf(f)

	bf.Run()

	got := bf.DataCells[:7]

	want := []int{0, 0, 72, 100, 87, 33, 10}

	for i, v := range got {
		if want[i] != v {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	}
}
