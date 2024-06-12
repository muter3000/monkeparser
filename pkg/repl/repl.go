package repl

import (
	"bufio"
	"fmt"
	"github.com/muter3000/monkeparser/pkg/lexer"
	"github.com/muter3000/monkeparser/pkg/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, err := r.output.Write(
				[]byte(fmt.Sprintf("%+v\n", tok)))
			if err != nil {
				return
			}
		}
	}
}
