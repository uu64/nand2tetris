package codewriter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CodeWriter struct {
	writer *bufio.Writer
}

func New(f io.Writer) *CodeWriter {
	return &CodeWriter{
		writer: bufio.NewWriter(f),
	}
}

func (cw *CodeWriter) WritePush(seg SegmentType, index int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "push %s %d\n", seg, index)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WritePop(seg SegmentType, index int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "pop %s %d\n", seg, index)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WriteArithmetic() error {
	// TODO: impl
	return nil
}

func (cw *CodeWriter) WriteLabel(label string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "label %s\n", label)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WriteGoTo(label string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "goto %s\n", label)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WriteIf(label string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "if-goto %s\n", label)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WriteCall(name string, nArgs int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "call %s %d\n", name, nArgs)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WriteFunction(name string, nLocals int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "function %s %d\n", name, nLocals)

	_, err := cw.writer.WriteString(b.String())
	return err
}

func (cw *CodeWriter) WriteReturn() error {
	_, err := cw.writer.WriteString("return")
	return err
}

func (cw *CodeWriter) Close() error {
	if err := cw.writer.Flush(); err != nil {
		return err
	}
	return nil
}
