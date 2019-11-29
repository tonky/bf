package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func eval(s io.Reader) (string, error) {
	buf := make([]byte, 1)
	cell := 0

	for {
		c, err := s.Read(buf)

		fmt.Println("read: ", string(buf), c, err)

		if err != nil && err != io.EOF {
			return "", err
		}

		if err == io.EOF || c == 0 {
			break
		}

		op := string(buf)

		switch op {
		// case '>':
		// 	program = append(program, Instruction{op_inc_dp, 0})
		// case '<':
		// 	program = append(program, Instruction{op_dec_dp, 0})
		case "+":
			fmt.Println("DEBUG: adding 1 from ", cell)
			cell += 1
		case "-":
			fmt.Println("DEBUG: substracting 1 from ", cell)
			cell -= 1
		case ".":
			fmt.Printf("Output via '.': %d\n", cell)
		}
	}

	return strconv.Itoa(cell), nil
}

func main() {
	// reader := bufio.NewReader(os.Stdin)

	// read_val, _ := reader.ReadByte()

	got, err := eval(bytes.NewBufferString("++."))

	fmt.Println("print eval: ", got, err)
}
