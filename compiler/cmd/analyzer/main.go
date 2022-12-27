package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/uu64/nand2tetris/compiler/internal/compile"
	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Cmd struct {
	source string
	output string
}

func (cmd *Cmd) write(b []byte) error {
	f, err := os.Create(cmd.output)
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

func (cmd *Cmd) compile() (*bytes.Buffer, error) {
	f, err := os.Open(cmd.source)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	compiler := compile.New(token.New(f))
	class, err := compiler.CompileClass()
	if err != nil {
		return nil, err
	}

	b, err := xml.MarshalIndent(class, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

func New(source, output string) *Cmd {
	return &Cmd{source, output}
}

func (cmd *Cmd) Run() (err error) {
	buf, err := cmd.compile()
	if err != nil {
		return err
	}

	err = cmd.write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Println("usage: jackc input")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		return
	}

	abs, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(abs)
	base := filepath.Base(abs)
	ext := filepath.Ext(abs)
	if ext != ".jack" {
		log.Fatalf("invalid extention: %s\n", ext)
	}
	output := fmt.Sprintf("%s/%s.xml", dir, base[0:len(base)-len(ext)])

	// TODO: directoryかファイル単体を渡す
	cmd := New(abs, output)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
