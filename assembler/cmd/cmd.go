package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/uu64/nand2tetris/assembler/internal/code"
	"github.com/uu64/nand2tetris/assembler/internal/parser"
)

type Cmd struct {
	asmfilePath string
}

func New(asmfilePath string) *Cmd {
	return &Cmd{asmfilePath}
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
	for p.HasMoreCommands() {
		p.Advance()
		switch p.CommandType() {
		case parser.A_CMD:
			s, err := strconv.Atoi(p.Symbol())
			if err != nil {
				return buf, err
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
	// parse
	buf, err := cmd.parse()
	if err != nil {
		return err
	}

	// write
	err = cmd.write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
