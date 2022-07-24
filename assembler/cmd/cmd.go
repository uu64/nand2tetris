package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/uu64/nand2tetris/assembler/internal/code"
	"github.com/uu64/nand2tetris/assembler/internal/parser"
	"github.com/uu64/nand2tetris/assembler/internal/symboltable"
)

type Cmd struct {
	asmfilePath string
	symbolTable *symboltable.SymbolTable
}

func New(asmfilePath string) *Cmd {
	return &Cmd{
		asmfilePath: asmfilePath,
		symbolTable: symboltable.New(),
	}
}

func (cmd *Cmd) scanSymbol() error {
	// create parser
	f, err := os.Open(cmd.asmfilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	p := parser.New(f)

	// scan
	romAddr := 0
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.A_CMD, parser.C_CMD:
			romAddr += 1
		case parser.L_CMD:
			cmd.symbolTable.AddEntry(p.Symbol(), romAddr)
		}
	}

	return nil
}

func (cmd *Cmd) parse() (*bytes.Buffer, error) {
	// create parser
	f, err := os.Open(cmd.asmfilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p := parser.New(f)

	// parse
	buf := bytes.NewBuffer([]byte{})
	ramAddr := 16
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.A_CMD:
			symbol := p.Symbol()
			addr, err := strconv.Atoi(symbol)
			if err != nil {
				if cmd.symbolTable.Contains(symbol) {
					addr = cmd.symbolTable.GetAddress(symbol)
				} else {
					cmd.symbolTable.AddEntry(symbol, ramAddr)
					addr = ramAddr
					ramAddr += 1
				}
			}
			fmt.Fprintf(buf, "%016b\n", addr)
		case parser.C_CMD:
			comp := code.Comp(p.Comp())
			dest := code.Dest(p.Dest())
			jump := code.Jump(p.Jump())
			fmt.Fprintf(buf, "111%07s%03s%03s\n", comp, dest, jump)
		case parser.L_CMD:
			// do nothing
		}
	}

	return buf, nil
}

func (cmd *Cmd) write(b []byte) error {
	outputPath := fmt.Sprintf("%s.hack", cmd.asmfilePath[0:len(cmd.asmfilePath)-len(".asm")])

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Cmd) Run() (err error) {
	if err = cmd.scanSymbol(); err != nil {
		return err
	}

	buf, err := cmd.parse()
	if err != nil {
		return err
	}

	err = cmd.write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
