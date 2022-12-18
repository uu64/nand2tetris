package compile

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Compiler struct {
	tokenizer *token.Tokenizer
}

func New(t *token.Tokenizer) *Compiler {
	return &Compiler{t}
}

func (c *Compiler) CompileClass() (*Class, error) {
	class := &Class{Tokens: []token.Token{}}

	if !c.tokenizer.HasMoreTokens() {
		return nil, fmt.Errorf("CompileClass: no tokens")
	}

	// class で始まるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}
	kwd, err := c.tokenizer.Keyword()
	if err != nil {
		return nil, err
	}
	if kwd.Val() != token.KwdClass {
		return nil, fmt.Errorf("CompileClass: class should start with CLASS")
	}
	class.Tokens = append(class.Tokens, *kwd)

	// className があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}
	className, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, err
	}
	class.Tokens = append(class.Tokens, *className)

	// { があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}
	open, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, err
	}
	if open.Val() != rune('{') {
		return nil, fmt.Errorf("CompileClass: symbol '{' is missing")
	}
	class.Tokens = append(class.Tokens, *open)

	// CompileClassVarDec()

	// CompileSubroutineVarDec()

	// } があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}
	close, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, err
	}
	if close.Val() != rune('}') {
		return nil, fmt.Errorf("CompileClass: symbol '}' is missing")
	}
	class.Tokens = append(class.Tokens, *close)

	return class, nil
}
