package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	compiler "github.com/uu64/nand2tetris/compiler/internal/engine"
	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
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

func (cmd *Cmd) encodeXML(class *compiler.Class) ([]byte, error) {
	return xml.MarshalIndent(class, "", "  ")
}

func (cmd *Cmd) compile() (*compiler.Class, error) {
	f, err := os.Open(cmd.source)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	compiler, err := compiler.New(tokenizer.New(f))
	if err != nil {
		return nil, err
	}

	return compiler.CompileClass()
}

func New(source, output string) *Cmd {
	return &Cmd{source, output}
}

func (cmd *Cmd) Run() (err error) {
	class, err := cmd.compile()
	if err != nil {
		return err
	}

	b, err := cmd.encodeXML(class)
	if err != nil {
		return err
	}

	err = cmd.write(b)
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
