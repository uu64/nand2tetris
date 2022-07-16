package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/uu64/nand2tetris/assembler/internal/code"
	"github.com/uu64/nand2tetris/assembler/internal/parser"
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
	ext := filepath.Ext(asmfile)
	hackfile := fmt.Sprintf("%s.hack", asmfile[0:len(asmfile)-len(ext)])

	// create parser
	in, err := os.Open(asmfile)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	p := parser.New(in)

	// parse
	buf := bytes.NewBuffer([]byte{})
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.A_CMD:
			s, err := strconv.Atoi(p.Symbol())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(buf, "%016b\n", s)
		case parser.C_CMD:
			comp := code.Comp(p.Comp())
			dest := code.Dest(p.Dest())
			jump := code.Jump(p.Jump())
			fmt.Fprintf(buf, "111%07s%03s%03s\n", comp, dest, jump)
		case parser.L_CMD:
		}
	}

	// write
	out, err := os.Create(hackfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
