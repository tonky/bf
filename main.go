package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Bf hold all interpreter state and methods
type Bf struct {
	inputReader        *bufio.Reader
	Input              []string // cached input
	InstructionPointer int

	DataCells      [30000]int
	DataPointer    int
	maxDataPointer int // pretty print helper to not pring while 'DataCells'
}

func newBf(input io.Reader) Bf {
	return Bf{
		inputReader:        bufio.NewReader(input),
		Input:              []string{},
		InstructionPointer: 0,
		DataCells:          [30000]int{},
		DataPointer:        0,
	}
}

// Loop is mutually recursive with 'Eval()', descending into in on each loop
func (bf *Bf) Loop() {
	startIdx := bf.InstructionPointer

	for {
		// Instruction is "[" - check if we need to skip this loop
		if bf.DataCells[bf.DataPointer] == 0 {
			bf.SKipToLoopEnd()
			return
		}

		bf.Eval()

		// Instruction is "]" - run loop again if data cell isn't 0
		if bf.DataCells[bf.DataPointer] > 0 {
			bf.InstructionPointer = startIdx
		}

		if bf.DataCells[bf.DataPointer] == 0 {
			return
		}
	}
}

// Eval is the main interpreter - evaluate and recurse into any loop we find
func (bf *Bf) Eval() {
	for {
		instruction, ok := bf.ReadAndAdvance()

		if !ok {
			return
		}

		switch instruction {
		case "+":
			bf.DataCells[bf.DataPointer] += 1
		case "-":
			bf.DataCells[bf.DataPointer] -= 1
		case "<":
			bf.DataPointer -= 1
		case ">":
			bf.DataPointer += 1
		case ".":
			fmt.Printf("%s", string(bf.DataCells[bf.DataPointer]))
		case ",":
			fmt.Printf("> Input: ")
			if c, err := bufio.NewReader(os.Stdin).ReadByte(); err != nil {
				panic("can't read from stdin")
			} else {
				bf.DataCells[bf.DataPointer] = int(c)
			}
		case "[":
			bf.Loop()
		case "]": // found end of loop into inner 'Eval()', just return to parent 'Loop'
			return
		case "j": // custom operator to increment DataCell pointer twice
			bf.DataPointer += 2
		default:
		}

		// save the last toched data cell to limit displaying DataCells output
		if bf.maxDataPointer < bf.DataPointer {
			bf.maxDataPointer = bf.DataPointer
		}
	}
}

// ReadAndAdvance reads current instruction and advances instruction pointer
func (bf *Bf) ReadAndAdvance() (string, bool) {
	ins, ok := bf.ReadInstruction()

	bf.InstructionPointer += 1

	return ins, ok
}

// ReadInstruction tries to return an instruction at 'InstructionPointer' postion,
// reading it from cahced Input if in bounds, or trying to read it from bufio.Reader
func (bf *Bf) ReadInstruction() (string, bool) {
	if bf.InstructionPointer < len(bf.Input) {
		return bf.Input[bf.InstructionPointer], true
	}

	if bf.InstructionPointer < 0 || bf.InstructionPointer > len(bf.Input) {
		fmt.Printf("IP: %d, Input: %d\n", bf.InstructionPointer, len(bf.Input))
		panic("Instruction pointer advanced beyond known data")
	}

	// we don't have cached Input instruction,so let's try reading it from Reader
	if c, err := bf.inputReader.ReadByte(); err != nil {
		return "", false
	} else {
		bf.Input = append(bf.Input, string(c))
	}

	return bf.Input[bf.InstructionPointer], true
}

// IntString pritns DataCells content as human readable string representation
func (bf Bf) IntString() string {
	res := ""

	for _, v := range bf.DataCells[:bf.maxDataPointer+1] {
		res += strconv.Itoa(v)
	}

	return res
}

// SkipToLoopEnd tries to advance to a matching "]" bracket in the input
func (bf *Bf) SKipToLoopEnd() {
	for matchingBracket := 1; matchingBracket > 0; bf.InstructionPointer++ {
		ins, ok := bf.ReadInstruction()

		if !ok {
			return
		}

		if ins == "[" {
			matchingBracket += 1
		}

		if ins == "]" {
			matchingBracket -= 1
		}
	}
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
		bf.Eval()

		return
	} else {
		fmt.Println("Can't find provided argument as file.")
	}

	fmt.Printf("Running program as a source string: '%s'\n", args[0])

	bf := newBf(bytes.NewBufferString(args[0]))

	bf.Eval()
}
