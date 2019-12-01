# Brainf*ck interpreter in Go

### Requirements

Prepare a brainf*ck interpreter

- Make it stack-based. This should mean that every loop must be in stack one level below current.

- Make it read input from io.Reader without knowing all input at once(no ioutil.ReadAll and similar, only read what is needed for current operation)

- Make it extensible. Let it have ability to add custom operations (for bonus points)

### Running

    $ go run main.go hello_world.bf
    Checking file existence: 'hello_world.bf'
    Exists! Opening now...
    Running program from file...
    Hello World!

    $ go run main.go '++>+++++[<+>-]++++++++[<++++++>-]<.'
    Checking file existence: '++>+++++[<+>-]++++++++[<++++++>-]<.'
    Can't find provided argument as file.
    Running program as a source string: '++>+++++[<+>-]++++++++[<++++++>-]<.'
    7%

### Testing

    $ go test
    7Hello World!
    PASS
    ok  github.com/tonky/bf 0.003s
