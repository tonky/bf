package main

import (
	"bytes"
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
	// hw := `++++++++ [ >++++  [ >++ >+++ >+++ >+  <<<<- ] >+ >+ >-  >>+ [<] <- ] `
	// hw2 := `++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++.`
	hw3 := `[ This program prints "Hello World!" and a newline to the screen, its
  length is 106 active command characters. [It is not the shortest.]

  This loop is an "initial comment loop", a simple way of adding a comment
  to a BF program such that you don't have to worry about any command
  characters. Any ".", ",", "+", "-", "<" and ">" characters are simply
  ignored, the "[" and "]" characters just have to be balanced. This
  loop and the commands it contains are ignored because the current cell
  defaults to a value of 0; the 0 value causes this loop to be skipped.
]
++++++++               Set Cell #0 to 8
[
    >++++               Add 4 to Cell #1; this will always set Cell #1 to 4
    [                   as the cell will be cleared by the loop
        >++             Add 2 to Cell #2
        >+++            Add 3 to Cell #3
        >+++            Add 3 to Cell #4
        >+              Add 1 to Cell #5
        <<<<-           Decrement the loop counter in Cell #1
    ]                   Loop till Cell #1 is zero; number of iterations is 4
    >+                  Add 1 to Cell #2
    >+                  Add 1 to Cell #3
    >-                  Subtract 1 from Cell #4
    >>+                 Add 1 to Cell #6
    [<]                 Move back to the first zero cell you find; this will
                        be Cell #1 which was cleared by the previous loop
    <-                  Decrement the loop Counter in Cell #0
]   `

	bf := newBf(bytes.NewBufferString(hw3))

	bf.Run()

	got := bf.DataCells[:7]

	want := []int{0, 0, 72, 104, 88, 32, 8, 0}

	for i, v := range got {
		if want[i] != v {
			t.Errorf("want: '%v', got: '%v'", want, got)
		}
	}
}
