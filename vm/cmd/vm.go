package cmd

import (
	"os"

	"github.com/uu64/nand2tetris/vm/internal/codewriter"
	"github.com/uu64/nand2tetris/vm/internal/parser"
)

type Cmd struct {
	vmfilePath  string
	asmfilePath string
}

func New(vmfilePath, asmfilePath string) *Cmd {
	return &Cmd{
		vmfilePath:  vmfilePath,
		asmfilePath: asmfilePath,
	}
}

func (cmd *Cmd) Run() (err error) {
	in, err := os.Open(cmd.vmfilePath)
	if err != nil {
		return err
	}
	defer in.Close()
	p := parser.New(in)

	out, err := os.Create(cmd.asmfilePath)
	if err != nil {
		return err
	}
	defer out.Close()
	cw := codewriter.New(out)

	// TODO: support multiple files
	cw.SetFileName(cmd.vmfilePath)

	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			return err
		}

		switch p.CommandType() {
		case parser.C_ARITHMETRIC:
			cw.WriteArithmetic(p.Arg1())
		case parser.C_PUSH:
			cw.WritePushPop(parser.C_PUSH, p.Arg1(), p.Arg2())
		case parser.C_POP:
			cw.WritePushPop(parser.C_POP, p.Arg1(), p.Arg2())
		default:
			// do nothing
		}
	}

	if err := cw.Close(); err != nil {
		return err
	}
	return nil
}
