package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/uu64/nand2tetris/compiler/internal/engine"
	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

const (
	extJack = ".jack"
)

type Cmd struct {
	source    string
	xmlOutput string
	vmOutput  string
}

func (cmd *Cmd) write(b []byte) error {
	f, err := os.Create(cmd.xmlOutput)
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

func (cmd *Cmd) encodeXML(class *engine.Class) ([]byte, error) {
	return xml.MarshalIndent(class, "", "  ")
}

func (cmd *Cmd) compile() (*engine.Class, error) {
	src, err := os.Open(cmd.source)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(cmd.vmOutput)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	compiler, err := engine.New(tokenizer.New(src), out)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Compiling %s\n", cmd.source)
	return compiler.CompileClass()
}

func New(source, xmlOutput, vmOutput string) *Cmd {
	return &Cmd{source, xmlOutput, vmOutput}
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

	filepath.Walk(abs, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), extJack) {
			xmlOutput := fmt.Sprintf("%s.xml", path[0:len(path)-len(extJack)])
			vmOutput := fmt.Sprintf("%s.vm", path[0:len(path)-len(extJack)])

			cmd := New(path, xmlOutput, vmOutput)
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})
}
