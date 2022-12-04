package cmd

import (
	"fmt"
	"os"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Cmd struct {
	source string
}

func New(source string) *Cmd {
	return &Cmd{source}
}

func (cmd *Cmd) Run() (err error) {
	f, err := os.Open(cmd.source)
	if err != nil {
		return err
	}
	defer f.Close()

	t := token.New(f)
	for t.HasMoreTokens() {
		if err := t.Advance(); err != nil {
			return err
		}

		tkType := t.TokenType()
		switch tkType {
		case token.TkKeyword:
			if v, err := t.Keyword(); err != nil {
				return err
			} else {
				fmt.Printf("kwd: %s\n", v)
			}
		case token.TkSymbol:
			if v, err := t.Symbol(); err != nil {
				return err
			} else {
				fmt.Printf("symbol: %s\n", v)
			}
		case token.TkIdentifier:
			if v, err := t.Identifier(); err != nil {
				return err
			} else {
				fmt.Printf("id: %s\n", v)
			}
		case token.TkIntConst:
			if v, err := t.IntVal(); err != nil {
				return err
			} else {
				fmt.Printf("int: %d\n", v)
			}
		case token.TkStringConst:
			if v, err := t.StringVal(); err != nil {
				return err
			} else {
				fmt.Printf("str: %s\n", v)
			}
		}
	}

	return nil
}
