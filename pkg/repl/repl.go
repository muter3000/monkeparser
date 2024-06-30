package repl

import (
	"bufio"
	"fmt"
	"github.com/muter3000/monkeparser/pkg/evaluator"
	"github.com/muter3000/monkeparser/pkg/lexer"
	"github.com/muter3000/monkeparser/pkg/object"
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

func (r *Repl) Start() {
	scanner := bufio.NewScanner(r.input)
	env := object.NewEnvironment()
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
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			_, err := io.WriteString(r.output, evaluated.Inspect())
			if err != nil {
				panic(err)
			}
			_, err = io.WriteString(r.output, "\n")
			if err != nil {
				return
			}
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
