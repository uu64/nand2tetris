package codewriter

import (
	"bufio"
	"fmt"
	"io"

	"github.com/uu64/nand2tetris/vm/internal/codewriter/code"
	"github.com/uu64/nand2tetris/vm/internal/parser"
)

type CodeWriter struct {
	writer        *bufio.Writer
	inputFileName string
}

func New(f io.Writer) *CodeWriter {
	return &CodeWriter{bufio.NewWriter(f), ""}
}

func (cw *CodeWriter) SetFileName(name string) {
	cw.inputFileName = name
	// TODO: impl something
}

func (cw *CodeWriter) Close() error {
	if err := cw.writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteArithmetic(cmd string) {
	switch cmd {
	case "add":
		cw.writer.WriteString(code.Add)
	default:
		// do notching
	}
}

func (cw *CodeWriter) WritePushPop(cmd parser.Cmd, segment string, index int) {
	switch cmd {
	case parser.C_PUSH:
		cw.writePush(cmd, segment, index)
	case parser.C_POP:
		cw.writePop(cmd, segment, index)
	default:
		// do nothing
	}
}

func (cw *CodeWriter) writePush(cmd parser.Cmd, segment string, index int) {
	switch segment {
	case "constant":
		cw.writer.WriteString(fmt.Sprintf(code.PushConstant, index))
	default:
		// do nothing
	}
}

func (cw *CodeWriter) writePop(cmd parser.Cmd, segment string, index int) {
	switch segment {
	case "constant":
		cw.writer.WriteString(code.PopConstant)
	default:
		// do nothing
	}
}
