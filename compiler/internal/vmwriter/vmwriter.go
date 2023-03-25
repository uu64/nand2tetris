package vmwriter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type VMWriter struct {
	writer *bufio.Writer
}

func New(f io.Writer) *VMWriter {
	return &VMWriter{
		writer: bufio.NewWriter(f),
	}
}

func (vw *VMWriter) WritePush(seg SegmentType, index int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "push %s %d\n", seg, index)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WritePop(seg SegmentType, index int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "pop %s %d\n", seg, index)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteArithmetic(cmd CommandType) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", cmd)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteLabel(label string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "label %s\n", label)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteGoTo(label string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "goto %s\n", label)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteIf(label string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "if-goto %s\n", label)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteCall(name string, nArgs int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "call %s %d\n", name, nArgs)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteFunction(name string, nLocals int) error {
	var b strings.Builder
	fmt.Fprintf(&b, "function %s %d\n", name, nLocals)

	_, err := vw.writer.WriteString(b.String())
	return err
}

func (vw *VMWriter) WriteReturn() error {
	_, err := vw.writer.WriteString("return")
	return err
}

func (vw *VMWriter) Close() error {
	if err := vw.writer.Flush(); err != nil {
		return err
	}
	return nil
}
