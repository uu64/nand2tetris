package codewriter

import (
	"bufio"
	"bytes"
	"fmt"
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

func (cw *CodeWriter) WriteArithmetic(cmd string) error {
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
		return fmt.Errorf("undefined operator: %s", cmd)
	}
	return nil
}

func (cw *CodeWriter) WritePushPop(cmd parser.Cmd, segment string, index int) error {
	switch cmd {
	case parser.C_PUSH:
		return cw.writePush(cmd, segment, index)
	case parser.C_POP:
		return cw.writePop(cmd, segment, index)
	default:
		return fmt.Errorf("invalid operation: %d", cmd)
	}
}

var memSegMap = map[string]string{
	parser.SEG_LOCAL: "LCL",
	parser.SEG_ARG:   "ARG",
	parser.SEG_THIS:  "THIS",
	parser.SEG_THAT:  "THAT",
	// base address
	parser.SEG_PTR:  "3",
	parser.SEG_TEMP: "5",
}

func (cw *CodeWriter) writePush(cmd parser.Cmd, segment string, index int) error {
	var b bytes.Buffer

	switch segment {
	case parser.SEG_CONST:
		b.WriteString(fmt.Sprintf("@%d\n", index))
		b.WriteString("D=A\n")
	case parser.SEG_LOCAL, parser.SEG_ARG, parser.SEG_THIS, parser.SEG_THAT:
		b.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		b.WriteString("D=M\n")
		b.WriteString(fmt.Sprintf("@%d\n", index))
		b.WriteString("A=D+A\n")
		b.WriteString("D=M\n")
	case parser.SEG_PTR, parser.SEG_TEMP:
		b.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		b.WriteString("D=A\n")
		b.WriteString(fmt.Sprintf("@%d\n", index))
		b.WriteString("A=D+A\n")
		b.WriteString("D=M\n")
	default:
		return fmt.Errorf("undefined segment: %s", segment)
	}

	b.WriteString("@SP\n")
	b.WriteString("A=M\n")
	b.WriteString("M=D\n")
	b.WriteString("@SP\n")
	b.WriteString("M=M+1\n")

	cw.writer.WriteString(b.String())
	return nil
}

func (cw *CodeWriter) writePop(cmd parser.Cmd, segment string, index int) error {
	var b bytes.Buffer

	if segment == parser.SEG_CONST {
		b.WriteString("@SP\n")
		b.WriteString("AM=M-1\n")
		cw.writer.WriteString(b.String())
		return nil
	}

	switch segment {
	case parser.SEG_LOCAL, parser.SEG_ARG, parser.SEG_THIS, parser.SEG_THAT:
		// calculate address
		b.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		b.WriteString("D=M\n")
		b.WriteString(fmt.Sprintf("@%d\n", index))
		b.WriteString("D=D+A\n")
		// save address
		b.WriteString("@R13\n")
		b.WriteString("M=D\n")
	case parser.SEG_PTR, parser.SEG_TEMP:
		// calculate address
		b.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		b.WriteString("D=A\n")
		b.WriteString(fmt.Sprintf("@%d\n", index))
		b.WriteString("D=D+A\n")
		// save address
		b.WriteString("@R13\n")
		b.WriteString("M=D\n")
	default:
		return fmt.Errorf("undefined segment: %s", segment)
	}

	b.WriteString("@SP\n")
	b.WriteString("AM=M-1\n")
	b.WriteString("D=M\n")

	b.WriteString("@R13\n")
	b.WriteString("A=M\n")
	b.WriteString("M=D\n")

	cw.writer.WriteString(b.String())
	return nil
}
