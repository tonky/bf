package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

// read RawInput byte by byte
// add to Input
// process instruction
// advance insruction pointer
// loops
type Bf struct {
	rawInput io.Reader

	Input              []string
	InstructionPointer int

	DataCells      *[30]int
	DataPointer    int
	maxDataPointer int
}

func (bf Bf) String() string {
	return fmt.Sprintf("==============\nBF struct\nInput: %s\nInstructionPointer: %d\nDataCells: %v\nDataPointer: %d\n=============",
		bf.Input, bf.InstructionPointer, bf.DataCells, bf.DataPointer)
}

func newBf(input io.Reader) Bf {
	var dc [30]int

	return Bf{
		rawInput:           input,
		Input:              []string{},
		InstructionPointer: 0,
		DataCells:          &dc,
		DataPointer:        0,
	}
}

// eiter read from Reader and add it to Input, or return data from Input
func (bf *Bf) ReadInstruction() (bool, string) {
	// fmt.Printf("> trying to read %d @ %v", bf.InstructionPointer, bf.Input)

	// if bf.InstructionPointer == len(bf.Input) {
	if bf.InstructionPointer < len(bf.Input) {
		// fmt.Println("DEBUG: Read from cached input")
		return true, bf.Input[bf.InstructionPointer]
	} else {
		// fmt.Println("DEBUG: Can't ead from cached input")
	}

	// fmt.Println("DEBUG: Trying to read from Reader")

	if ok, c := readNextChar(bf.rawInput); ok {
		// fmt.Println("> read input byte: ", c)

		bf.Input = append(bf.Input, c)
		// fmt.Println("> enhanced input, now: ", bf.Input)
	} else {
		// fmt.Println("> no more reading from input: ", c)
		return false, ""
	}

	if bf.InstructionPointer < 0 || bf.InstructionPointer > len(bf.Input) {
		panic("undefined behavior: can't point behind known chars")
	}

	// fmt.Printf("> ReadInstruction: pos '%d': '%s', data: %d @ %d\n", bf.InstructionPointer, bf.Input[bf.InstructionPointer], bf.DataCells[bf.DataPointer], bf.DataPointer)

	return true, bf.Input[bf.InstructionPointer]
}

func (bf *Bf) IncrementCurrentCell() {
	if bf.DataCells[bf.DataPointer] >= 255 {
		// fmt.Println("data cell value overflow:")
		panic(bf)
	}

	bf.DataCells[bf.DataPointer] += 1

	// fmt.Printf("Inc: %d @ %d\n", bf.DataCells[bf.DataPointer], bf.DataPointer)
}

func (bf *Bf) DecrementCurrentCell() {
	if bf.DataCells[bf.DataPointer] < 0 {
		// fmt.Println("Trying to decrement cell under zero!")
		panic(bf)
	}

	bf.DataCells[bf.DataPointer] -= 1
}

func (bf *Bf) IncrementDataPointer() {
	bf.DataPointer += 1

	if bf.DataPointer > bf.maxDataPointer {
		bf.maxDataPointer = bf.DataPointer
	}
}

func (bf *Bf) DecrementDataPointer() {
	if bf.DataPointer <= 0 {
		// fmt.Println("Trying to decrement data pointer below zero!")
		panic(bf)
	}

	bf.DataPointer -= 1
}

func (bf *Bf) AdvanceInstructionPointer() {
	bf.InstructionPointer += 1
}

func (bf Bf) PrintCurrentCell() {
	fmt.Printf("%s", string(bf.DataCells[bf.DataPointer]))
}

func (bf *Bf) SkipLoop() {
	bracketBalance := 0

	for {
		ok, ins := bf.ReadInstruction()

		if !ok {
			return
		}

		if ins == "[" {
			bracketBalance += 1
		}

		if ins == "]" {
			bracketBalance -= 1
		}

		bf.AdvanceInstructionPointer()

		if bracketBalance == 0 {
			return
		}
	}
}

func (bf *Bf) Run() {
	for {
		ok, ins := bf.ReadInstruction()

		if !ok {
			return
		}

		switch ins {
		case "[":
			bf.Loop()
		case "]": // we're inside loop and it's the end of it
			// fmt.Println("end of loop in inner, returning to parent")
			// bf.AdvanceInstructionPointer()
			return
		case "+":
			bf.IncrementCurrentCell()
			bf.AdvanceInstructionPointer()
		case "-":
			bf.DecrementCurrentCell()
			bf.AdvanceInstructionPointer()
		case ">":
			bf.IncrementDataPointer()
			bf.AdvanceInstructionPointer()
		case "<":
			bf.DecrementDataPointer()
			bf.AdvanceInstructionPointer()
		case ".":
			bf.PrintCurrentCell()
			bf.AdvanceInstructionPointer()
		default: // skip unknown commands
			// fmt.Printf("> unknown instruction: '%s', skipping", ins)
			bf.AdvanceInstructionPointer()
		}
	}
}

func (bf *Bf) Loop() {
	startIdx := bf.InstructionPointer

	// fmt.Println("> Starting 'Loop', startIdx: ", startIdx)

	for {
		ok, ins := bf.ReadInstruction()

		if !ok {
			return
		}

		if ins == "[" && bf.DataCells[bf.DataPointer] == 0 {
			// fmt.Println("> skipping loop from beginning")
			bf.SkipLoop()
			// fmt.Println("> returning from SkipLoop() to parent Run()")
			return
		}

		// fmt.Println(bf)

		if ins == "]" && bf.DataCells[bf.DataPointer] == 0 {
			// fmt.Println("> skipping loop end")
			bf.AdvanceInstructionPointer()
			return
		}

		if ins == "]" && bf.DataCells[bf.DataPointer] != 0 {
			bf.InstructionPointer = startIdx
			// fmt.Println("> ']' rewound loop to next position after: ", startIdx)
		}

		bf.AdvanceInstructionPointer()
		// fmt.Println("> Evaluating inner data, doing 'bf.Run()' from bf.Loop()")
		// fmt.Println(bf)
		bf.Run()
		// fmt.Printf("we're back in loop with starting index '%d'\n", startIdx)
	}
}

func (bf Bf) IntString() string {
	res := ""

	for _, v := range bf.DataCells[:bf.maxDataPointer+1] {
		res += strconv.Itoa(v)
	}

	return res
}

func (bf Bf) Ascii() string {
	res := ""

	for _, v := range bf.DataCells[:bf.maxDataPointer+1] {
		res += fmt.Sprint(string(v))
	}

	return res
}

func readNextChar(r io.Reader) (bool, string) {
	buf := make([]byte, 1)

	c, err := r.Read(buf)

	if err != nil && err != io.EOF {
		return false, ""
	}

	if err == io.EOF || c == 0 {
		return false, ""
	}

	return true, string(buf)
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Please provide an argument in form of filename or code")
		os.Exit(1)
	}

	// try checking if it's a file
	fmt.Printf("Checking file existence: '%s'\n", args[0])

	if _, err := os.Stat(args[0]); err == nil {
		fmt.Println("Exists! Opening now...")

		f, err := os.Open(args[0])

		if err != nil {
			panic(err)
		}

		fmt.Println("Running program from file...")

		bf := newBf(f)
		bf.Run()

		return
	} else {
		fmt.Println("Can't find provided argument as file.")
	}

	fmt.Printf("Running program as a source string: '%s'\n", args[0])

	bf := newBf(bytes.NewBufferString(args[0]))
	bf.Run()
}
