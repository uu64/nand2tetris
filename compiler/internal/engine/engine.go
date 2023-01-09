package engine

import (
	"fmt"

	token "github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

type Compiler struct {
	tokenizer *token.Tokenizer
}

func New(t *token.Tokenizer) (*Compiler, error) {
	if err := t.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}

	return &Compiler{t}, nil
}

func (c *Compiler) consumeKeyword(expected ...token.KeywordType) (*token.Keyword, error) {
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

func (c *Compiler) consumeSymbol(expected ...rune) (*token.Symbol, error) {
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

func (c *Compiler) consumeIdentifier() (*token.Identifier, error) {
	id, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("consumeIdentifier: %w", err)
	}

	return id, c.tokenizer.Advance()
}
