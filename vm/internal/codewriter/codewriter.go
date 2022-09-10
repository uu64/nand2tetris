package codewriter

import (
	"bufio"
	"io"

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
		cw.writer.WriteString(add())
	case "sub":
		cw.writer.WriteString(sub())
	case "neg":
		cw.writer.WriteString(neg())
	case "eq":
		cw.writer.WriteString(eq())
	case "gt":
		cw.writer.WriteString(gt())
	case "lt":
		cw.writer.WriteString(lt())
	case "and":
		cw.writer.WriteString(and())
	case "or":
		cw.writer.WriteString(or())
	case "not":
		cw.writer.WriteString(not())
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
		cw.writer.WriteString(pushConstant(index))
	default:
		// do nothing
	}
}

func (cw *CodeWriter) writePop(cmd parser.Cmd, segment string, index int) {
	switch segment {
	case "constant":
		cw.writer.WriteString(popConstant())
	default:
		// do nothing
	}
}
