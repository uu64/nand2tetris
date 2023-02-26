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

	token "github.com/uu64/nand2tetris/compiler/internal/tokenizer"
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

	tokens := &token.Tokens{
		Tokens: []token.Token{},
	}

	for t.HasMoreTokens() {
		if err := t.Advance(); err != nil {
			return nil, err
		}

		tkType := t.TokenType()
		switch tkType {
		case token.TkKeyword:
			v, err := t.Keyword()
			if err != nil {
				return nil, err
			}
			tokens.Tokens = append(tokens.Tokens, v)
		case token.TkSymbol:
			v, err := t.Symbol()
			if err != nil {
				return nil, err
			}
			tokens.Tokens = append(tokens.Tokens, v)
		case token.TkIdentifier:
			v, err := t.Identifier()
			if err != nil {
				return nil, err
			}
			tokens.Tokens = append(tokens.Tokens, v)
		case token.TkIntConst:
			v, err := t.IntVal()
			if err != nil {
				return nil, err
			}
			tokens.Tokens = append(tokens.Tokens, v)
		case token.TkStringConst:
			v, err := t.StringVal()
			if err != nil {
				return nil, err
			}
			tokens.Tokens = append(tokens.Tokens, v)
		default:
			continue
		}
	}

	b, err := xml.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
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
