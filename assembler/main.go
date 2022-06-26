package main

import (
	"log"
	"os"

	"github.com/uu64/nand2tetris/assembler/lib/parser"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	p := parser.New(f)

	for p.HasMoreCommands() {
		p.Advance()
		// fmt.Printf("%s\n", p.Symbol())
	}
}
