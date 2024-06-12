package main

import (
	"monkeparser/pkg/repl"
	"os"
)

func main() {
	rpl := repl.New(os.Stdin, os.Stdout, ">> ")
	rpl.Start()
}
