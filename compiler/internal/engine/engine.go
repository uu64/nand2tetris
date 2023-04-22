package engine

import (
	"fmt"
	"os"

	"github.com/uu64/nand2tetris/compiler/internal/symtab"
	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
	"github.com/uu64/nand2tetris/compiler/internal/vmwriter"
)

type Compiler struct {
	tokenizer  *tokenizer.Tokenizer
	symtab     *symtab.Symtab
	codewriter *vmwriter.VMWriter
}

func New(t *tokenizer.Tokenizer) (*Compiler, error) {
	if err := t.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}

	return &Compiler{
		tokenizer: t,
		symtab:    symtab.New(),
		// TODO: impl
		codewriter: vmwriter.New(os.Stdout),
		// codewriter: vmwriter.New(io.Discard),
	}, nil
}

func (c *Compiler) consumeKeyword(expected ...tokenizer.KeywordType) (*tokenizer.Keyword, error) {
	kwd, err := c.tokenizer.Keyword()
	if err != nil {
		return nil, fmt.Errorf("consumeKeyword: %w", err)
	}

	if len(expected) == 0 {
		return kwd, c.tokenizer.Advance()
	}

	for _, v := range expected {
		if kwd.Val() == v {
			return kwd, c.tokenizer.Advance()
		}
	}

	return nil, fmt.Errorf("consumeKeyword: expected %v, got %v", expected, kwd)
}

func (c *Compiler) consumeSymbol(expected ...rune) (*tokenizer.Symbol, error) {
	symbol, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("consumeSymbol: %w", err)
	}

	if len(expected) == 0 {
		return symbol, c.tokenizer.Advance()
	}

	for _, v := range expected {
		if symbol.Val() == v {
			return symbol, c.tokenizer.Advance()
		}
	}

	return nil, fmt.Errorf("consumeSymbol: expected %v, got %v", expected, symbol)
}
