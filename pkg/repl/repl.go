package repl

import (
	"bufio"
	"fmt"
	"github.com/muter3000/monkeparser/pkg/lexer"
	"github.com/muter3000/monkeparser/pkg/parser"
	"io"
)

type Repl struct {
	input  io.Reader
	output io.Writer

	prompt string
}

func New(input io.Reader, output io.Writer, prompt string) *Repl {
	return &Repl{input: input, output: output, prompt: prompt}
}

const welcomeMessage = `
Welcome to the Monke programming language REPL!
Feel free to type in commands (they won't work).
`

func (r *Repl) Start() {
	_, err := r.output.Write([]byte(welcomeMessage))
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(r.input)
	for {
		fmt.Printf(r.prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(r.output, p.Errors())
			continue
		}
		_, err := io.WriteString(r.output, program.String())
		if err != nil {
			panic(err)
		}
		_, err = io.WriteString(r.output, "\n")
		if err != nil {
			panic(err)
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, err := io.WriteString(out, "You wrote some really bad code!\n")
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(out, " parser errors:\n")
	if err != nil {
		panic(err)
	}
	for _, msg := range errors {
		_, err = io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			panic(err)
		}
	}
}
