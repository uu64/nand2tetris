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
	case "add", "sub", "and", "or":
		cw.binary(cmd)
	case "neg", "not":
		cw.unary(cmd)
	case "eq":
		cw.writer.WriteString(eq())
	case "gt":
		cw.writer.WriteString(gt())
	case "lt":
		cw.writer.WriteString(lt())
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}
	return nil
}

func (cw *CodeWriter) unary(cmd string) error {
	var b bytes.Buffer
	// this code is same as the code to pop from a constant segment
	b.WriteString("@SP\n")
	b.WriteString("AM=M-1\n")

	// update stack with the result
	switch cmd {
	case "neg":
		b.WriteString("M=-M\n")
	case "not":
		b.WriteString("M=!M\n")
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}

	// add stack pointer
	b.WriteString("@SP\n")
	b.WriteString("M=M+1\n")

	cw.writer.WriteString(b.String())
	return nil
}

func (cw *CodeWriter) binary(cmd string) error {
	var b bytes.Buffer
	// this code is same as the code to pop from a constant segment
	b.WriteString("@SP\n")
	b.WriteString("AM=M-1\n")

	// save the value and pop an another value
	b.WriteString("D=M\n")
	b.WriteString("@SP\n")
	b.WriteString("AM=M-1\n")

	// update stack with the result
	switch cmd {
	case "add":
		b.WriteString("M=M+D\n")
	case "sub":
		b.WriteString("M=M-D\n")
	case "and":
		b.WriteString("M=M&D\n")
	case "or":
		b.WriteString("M=M|D\n")
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}

	// add stack pointer
	b.WriteString("@SP\n")
	b.WriteString("M=M+1\n")

	cw.writer.WriteString(b.String())
	return nil
}

// func (cw *CodeWriter) cond(cmd string) error {}

func (cw *CodeWriter) WritePushPop(cmd parser.Cmd, segment string, index int) error {
	switch cmd {
	case parser.C_PUSH:
		return cw.writePush(segment, index)
	case parser.C_POP:
		return cw.writePop(segment, index)
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

// writePush outputs the asm code to push a value to a specific segment.
func (cw *CodeWriter) writePush(segment string, index int) error {
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
	case parser.SEG_STATIC:
		b.WriteString(fmt.Sprintf("@Xxx.%d\n", index))
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

// writePop outputs the asm code to pop a value from a specific segment.
// The return value is set to M.
func (cw *CodeWriter) writePop(segment string, index int) error {
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
	case parser.SEG_STATIC:
		b.WriteString(fmt.Sprintf("@Xxx.%d\n", index))
		b.WriteString("D=A\n")
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
