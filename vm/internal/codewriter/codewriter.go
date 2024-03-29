package codewriter

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/uu64/nand2tetris/vm/internal/parser"
)

type CodeWriter struct {
	writer        *bufio.Writer
	inputFileName string
	functionName  string
	counter       int
}

func New(f io.Writer) *CodeWriter {
	return &CodeWriter{
		writer:        bufio.NewWriter(f),
		inputFileName: "",
		functionName:  "",
		counter:       0,
	}
}

func (cw *CodeWriter) WriteInit() {
	cw.writer.WriteString("@256\n")
	cw.writer.WriteString("D=A\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=D\n")
	cw.WriteCall("Sys.init", 0)
}

func (cw *CodeWriter) SetFileName(name string) {
	cw.inputFileName = name
}

func (cw *CodeWriter) Close() error {
	if err := cw.writer.Flush(); err != nil {
		return err
	}
	return nil
}

func (cw *CodeWriter) WriteArithmetic(cmd string) error {
	switch cmd {
	case parser.CMD_ADD, parser.CMD_SUB, parser.CMD_AND, parser.CMD_OR:
		cw.binary(cmd)
	case parser.CMD_NEG, parser.CMD_NOT:
		cw.unary(cmd)
	case parser.CMD_EQ, parser.CMD_GT, parser.CMD_LT:
		cw.cond(cmd)
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}
	return nil
}

func (cw *CodeWriter) unary(cmd string) error {
	// this code is same as the code to pop from a constant segment
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")

	// update stack with the result
	switch cmd {
	case parser.CMD_NEG:
		cw.writer.WriteString("M=-M\n")
	case parser.CMD_NOT:
		cw.writer.WriteString("M=!M\n")
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}

	// add stack pointer
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=M+1\n")

	return nil
}

func (cw *CodeWriter) binary(cmd string) error {
	// this code is same as the code to pop from a constant segment
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")

	// save the value and pop an another value
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")

	// update stack with the result
	switch cmd {
	case parser.CMD_ADD:
		cw.writer.WriteString("M=M+D\n")
	case parser.CMD_SUB:
		cw.writer.WriteString("M=M-D\n")
	case parser.CMD_AND:
		cw.writer.WriteString("M=M&D\n")
	case parser.CMD_OR:
		cw.writer.WriteString("M=M|D\n")
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}

	// add stack pointer
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=M+1\n")

	return nil
}

func (cw *CodeWriter) cond(cmd string) error {
	id := strings.ToUpper(cmd)

	// this code is same as the code to pop from a constant segment
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")

	// save the value and pop an another value
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")

	// compare
	cw.writer.WriteString("D=M-D\n")
	cw.writer.WriteString(fmt.Sprintf("@%s%d_T\n", id, cw.counter))
	switch cmd {
	case parser.CMD_EQ:
		cw.writer.WriteString("D;JEQ\n")
	case parser.CMD_GT:
		cw.writer.WriteString("D;JGT\n")
	case parser.CMD_LT:
		cw.writer.WriteString("D;JLT\n")
	default:
		return fmt.Errorf("undefined operator: %s", cmd)
	}
	cw.writer.WriteString(fmt.Sprintf("@%s%d_F\n", id, cw.counter))
	cw.writer.WriteString("0;JMP\n")

	// set true or false
	cw.writer.WriteString(fmt.Sprintf("(%s%d_T)\n", id, cw.counter))
	cw.writer.WriteString("D=-1\n")
	cw.writer.WriteString(fmt.Sprintf("@%s%d_END\n", id, cw.counter))
	cw.writer.WriteString("0;JMP\n")
	cw.writer.WriteString(fmt.Sprintf("(%s%d_F)\n", id, cw.counter))
	cw.writer.WriteString("D=0\n")
	cw.writer.WriteString(fmt.Sprintf("@%s%d_END\n", id, cw.counter))
	cw.writer.WriteString("0;JMP\n")

	// update stack with the result
	cw.writer.WriteString(fmt.Sprintf("(%s%d_END)\n", id, cw.counter))
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("A=M\n")
	cw.writer.WriteString("M=D\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=M+1\n")

	cw.counter += 1
	return nil
}

func (cw *CodeWriter) WritePushPop(cmd parser.CmdType, segment string, index int) error {
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
	parser.SEG_PTR:   "3", // base address
	parser.SEG_TEMP:  "5", // base address
}

// writePush outputs the asm code to push a value to a specific segment.
func (cw *CodeWriter) writePush(segment string, index int) error {
	switch segment {
	case parser.SEG_CONST:
		cw.writer.WriteString(fmt.Sprintf("@%d\n", index))
		cw.writer.WriteString("D=A\n")
	case parser.SEG_LOCAL, parser.SEG_ARG, parser.SEG_THIS, parser.SEG_THAT:
		cw.writer.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		cw.writer.WriteString("D=M\n")
		cw.writer.WriteString(fmt.Sprintf("@%d\n", index))
		cw.writer.WriteString("A=D+A\n")
		cw.writer.WriteString("D=M\n")
	case parser.SEG_PTR, parser.SEG_TEMP:
		cw.writer.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		cw.writer.WriteString("D=A\n")
		cw.writer.WriteString(fmt.Sprintf("@%d\n", index))
		cw.writer.WriteString("A=D+A\n")
		cw.writer.WriteString("D=M\n")
	case parser.SEG_STATIC:
		ns := strings.TrimSuffix(cw.inputFileName, filepath.Ext(cw.inputFileName))
		cw.writer.WriteString(fmt.Sprintf("@%s.%d\n", ns, index))
		cw.writer.WriteString("D=M\n")
	default:
		return fmt.Errorf("undefined segment: %s", segment)
	}

	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("A=M\n")
	cw.writer.WriteString("M=D\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=M+1\n")

	return nil
}

// writePop outputs the asm code to pop a value from a specific segment.
// The return value is set to M.
func (cw *CodeWriter) writePop(segment string, index int) error {
	if segment == parser.SEG_CONST {
		cw.writer.WriteString("@SP\n")
		cw.writer.WriteString("AM=M-1\n")
		return nil
	}

	// calculate address
	switch segment {
	case parser.SEG_LOCAL, parser.SEG_ARG, parser.SEG_THIS, parser.SEG_THAT:
		cw.writer.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		cw.writer.WriteString("D=M\n")
		cw.writer.WriteString(fmt.Sprintf("@%d\n", index))
		cw.writer.WriteString("D=D+A\n")
	case parser.SEG_PTR, parser.SEG_TEMP:
		cw.writer.WriteString(fmt.Sprintf("@%s\n", memSegMap[segment]))
		cw.writer.WriteString("D=A\n")
		cw.writer.WriteString(fmt.Sprintf("@%d\n", index))
		cw.writer.WriteString("D=D+A\n")
	case parser.SEG_STATIC:
		ns := strings.TrimSuffix(cw.inputFileName, filepath.Ext(cw.inputFileName))
		cw.writer.WriteString(fmt.Sprintf("@%s.%d\n", ns, index))
		cw.writer.WriteString("D=A\n")
	default:
		return fmt.Errorf("undefined segment: %s", segment)
	}

	// save address
	cw.writer.WriteString("@R13\n")
	cw.writer.WriteString("M=D\n")

	// update the stack and the segment
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString("@R13\n")
	cw.writer.WriteString("A=M\n")
	cw.writer.WriteString("M=D\n")

	return nil
}

func (cw *CodeWriter) WriteLabel(label string) error {
	cw.writer.WriteString(fmt.Sprintf("(%s)\n", label))
	return nil
}

func (cw *CodeWriter) WriteGoto(label string) error {
	cw.writer.WriteString(fmt.Sprintf("@%s\n", label))
	cw.writer.WriteString("0;JMP\n")
	return nil
}

func (cw *CodeWriter) WriteIf(label string) error {
	// this code is same as the code to pop from a constant segment
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("AM=M-1\n")
	cw.writer.WriteString("D=M\n")

	// if M=0, do nothing
	cw.writer.WriteString(fmt.Sprintf("@IF%d\n", cw.counter))
	cw.writer.WriteString("D;JEQ\n")

	// if M!=0, jump to label
	cw.writer.WriteString(fmt.Sprintf("@%s\n", label))
	cw.writer.WriteString("0;JMP\n")

	cw.writer.WriteString(fmt.Sprintf("(IF%d)\n", cw.counter))

	cw.counter += 1
	return nil
}

func (cw *CodeWriter) WriteCall(functionName string, numArgs int) error {
	retAddr := fmt.Sprintf("%d.RET_ADDR", cw.counter)

	pushAddr := func(addr string) {
		cw.writer.WriteString(fmt.Sprintf("@%s\n", addr))
		cw.writer.WriteString("D=M\n")
		cw.writer.WriteString("@SP\n")
		cw.writer.WriteString("A=M\n")
		cw.writer.WriteString("M=D\n")
		cw.writer.WriteString("@SP\n")
		cw.writer.WriteString("M=M+1\n")
	}

	// push return-addr
	cw.writer.WriteString(fmt.Sprintf("@%s\n", retAddr))
	cw.writer.WriteString("D=A\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("A=M\n")
	cw.writer.WriteString("M=D\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=M+1\n")

	// push LCL
	// push ARG
	// push THIS
	// push THAT
	pushAddr(memSegMap[parser.SEG_LOCAL])
	pushAddr(memSegMap[parser.SEG_ARG])
	pushAddr(memSegMap[parser.SEG_THIS])
	pushAddr(memSegMap[parser.SEG_THAT])

	// ARG = SP-n-5
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString(fmt.Sprintf("@%d\n", numArgs))
	cw.writer.WriteString("D=D-A\n")
	cw.writer.WriteString("@5\n")
	cw.writer.WriteString("D=D-A\n")
	cw.writer.WriteString("@ARG\n")
	cw.writer.WriteString("M=D\n")

	// LCL = SP
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString("@LCL\n")
	cw.writer.WriteString("M=D\n")

	// goto f
	cw.WriteGoto(functionName)

	// (return-address)
	cw.writer.WriteString(fmt.Sprintf("(%s)\n", retAddr))

	cw.counter += 1
	return nil
}

func (cw *CodeWriter) WriteReturn() error {
	// FRAME = LCL
	cw.writer.WriteString("@LCL\n")
	cw.writer.WriteString("D=M\n")
	// RET = *(FRAME-5)
	cw.writer.WriteString("@5\n")
	cw.writer.WriteString("A=D-A\n")
	cw.writer.WriteString("D=M\n")

	// save RET
	cw.writer.WriteString("@R13\n")
	cw.writer.WriteString("M=D\n")

	// *ARG = pop()
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("A=M-1\n")
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString("@ARG\n")
	cw.writer.WriteString("A=M\n")
	cw.writer.WriteString("M=D\n")

	// SP = ARG+1
	cw.writer.WriteString("@ARG\n")
	cw.writer.WriteString("D=M\n")
	cw.writer.WriteString("@SP\n")
	cw.writer.WriteString("M=D+1\n")

	for i, seg := range []string{parser.SEG_THAT, parser.SEG_THIS, parser.SEG_ARG, parser.SEG_LOCAL} {
		// seg = *(FRAME-i)
		cw.writer.WriteString("@LCL\n")
		cw.writer.WriteString("D=M\n")
		cw.writer.WriteString(fmt.Sprintf("@%d\n", i+1))
		cw.writer.WriteString("A=D-A\n")
		cw.writer.WriteString("D=M\n")
		cw.writer.WriteString(fmt.Sprintf("@%s\n", memSegMap[seg]))
		cw.writer.WriteString("M=D\n")
	}

	// goto RET
	cw.writer.WriteString("@R13\n")
	cw.writer.WriteString("A=M\n")
	cw.writer.WriteString("0;JMP\n")
	return nil
}

func (cw *CodeWriter) WriteFunction(functionName string, numLocals int) error {
	cw.writer.WriteString(fmt.Sprintf("(%s)\n", functionName))
	for i := 0; i < numLocals; i++ {
		cw.writePush(parser.SEG_CONST, 0)
	}
	return nil
}
