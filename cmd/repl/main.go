package main

import (
	"github.com/muter3000/monkeparser/pkg/repl"
	"os"
)

func main() {
	rpl := repl.New(os.Stdin, os.Stdout, ">> ")
	rpl.Start()
}
