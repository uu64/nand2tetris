package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/uu64/nand2tetris/assembler/lib/code"
	"github.com/uu64/nand2tetris/assembler/lib/parser"
)

func usage() {
	fmt.Println("usage: asmc /path/to/file.asm")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	asmfile := os.Args[1]
	// ext := filepath.Ext(asmfile)
	// hackfile = fmt.Sprintf("%s.hack", asmfile[0:len(asmfile)-len(ext)])

	f, err := os.Open(asmfile)
	if err != nil {
		log.Fatal(err)
	}

	p := parser.New(f)

	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.A_CMD:
			// symbol := p.Symbol()
			s, err := strconv.Atoi(p.Symbol())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%016b\n", s)
		case parser.C_CMD:
			comp := p.Comp()
			dest := p.Dest()
			jump := p.Jump()
			fmt.Printf("111%07s%03s%03s\n", code.Comp(comp), code.Dest(dest), code.Jump(jump))
		case parser.L_CMD:
		}
	}
}
