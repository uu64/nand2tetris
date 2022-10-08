package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/uu64/nand2tetris/vm/internal/codewriter"
	"github.com/uu64/nand2tetris/vm/internal/parser"
)

type Cmd struct {
	vmfilePaths []string
	asmfilePath string
}

func New(vmfilePaths []string, asmfilePath string) *Cmd {
	return &Cmd{
		vmfilePaths: vmfilePaths,
		asmfilePath: asmfilePath,
	}
}

func parse(cw *codewriter.CodeWriter, asmfilePath string) error {
	in, err := os.Open(asmfilePath)
	if err != nil {
		return err
	}
	defer in.Close()
	p := parser.New(in)

	cw.SetFileName(filepath.Base(asmfilePath))

	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			return err
		}

		var err error
		switch p.CommandType() {
		case parser.C_ARITHMETRIC:
			err = cw.WriteArithmetic(p.Arg1())
		case parser.C_PUSH:
			err = cw.WritePushPop(parser.C_PUSH, p.Arg1(), p.Arg2())
		case parser.C_POP:
			err = cw.WritePushPop(parser.C_POP, p.Arg1(), p.Arg2())
		case parser.C_LABEL:
			err = cw.WriteLabel(p.Arg1())
		case parser.C_GOTO:
			err = cw.WriteGoto(p.Arg1())
		case parser.C_IF:
			err = cw.WriteIf(p.Arg1())
		case parser.C_FUNCTION:
			err = cw.WriteFunction(p.Arg1(), p.Arg2())
		case parser.C_RETURN:
			err = cw.WriteReturn()
		case parser.COMMENT, parser.EMPTY:
			// do nothing
		default:
			err = fmt.Errorf("undefined command type: %d", p.CommandType())
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (cmd *Cmd) Run() (err error) {
	out, err := os.Create(cmd.asmfilePath)
	if err != nil {
		return err
	}
	defer out.Close()
	cw := codewriter.New(out)

	for _, vmfilePath := range cmd.vmfilePaths {
		if err := parse(cw, vmfilePath); err != nil {
			return err
		}
	}

	if err := cw.Close(); err != nil {
		return err
	}
	return nil
}
