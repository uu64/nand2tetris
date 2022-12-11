package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func (cmd *Cmd) parse() (*bytes.Buffer, error) {
	f, err := os.Open(cmd.source)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := token.New(f)
	buf := bytes.NewBuffer([]byte{})
	fmt.Fprintf(buf, "<tokens>\n")

	for t.HasMoreTokens() {
		if err := t.Advance(); err != nil {
			return nil, err
		}

		tkType := t.TokenType()
		switch tkType {
		case token.TkKeyword:
			if v, err := t.Keyword(); err != nil {
				return nil, err
			} else {
				fmt.Fprintf(buf, "<keyword> %s </keyword>\n", v)
			}
		case token.TkSymbol:
			if v, err := t.Symbol(); err != nil {
				return nil, err
			} else {
				fmt.Fprintf(buf, "<symbol> %s </symbol>\n", v)
			}
		case token.TkIdentifier:
			if v, err := t.Identifier(); err != nil {
				return nil, err
			} else {
				fmt.Fprintf(buf, "<identifier> %s </identifier>\n", v)
			}
		case token.TkIntConst:
			if v, err := t.IntVal(); err != nil {
				return nil, err
			} else {
				fmt.Fprintf(buf, "<integerConstant> %d </integerConstant>\n", v)
			}
		case token.TkStringConst:
			if v, err := t.StringVal(); err != nil {
				return nil, err
			} else {
				fmt.Fprintf(buf, "<stringConstant> %s </stringConstant>\n", v)
			}
		default:
			continue
		}
	}

	fmt.Fprintf(buf, "</tokens>\n")
	return buf, nil
}

func New(source, output string) *Cmd {
	return &Cmd{source, output}
}

func (cmd *Cmd) Run() (err error) {
	buf, err := cmd.parse()
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
	output := fmt.Sprintf("%s/%sT.xml", dir, base[0:len(base)-len(ext)])

	// TODO: directoryかファイル単体を渡す
	cmd := New(abs, output)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
