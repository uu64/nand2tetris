package compile

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Compiler struct {
	tokenizer *token.Tokenizer
	class     *Class
}

func New(t *token.Tokenizer) *Compiler {
	class := &Class{Tokens: []token.Token{}}
	return &Compiler{t, class}
}

func (c *Compiler) CompileClass() (*Class, error) {
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
	c.class.Tokens = append(c.class.Tokens, *kwd)

	// className があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}
	className, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, err
	}
	c.class.Tokens = append(c.class.Tokens, *className)

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
	c.class.Tokens = append(c.class.Tokens, *open)

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
	c.class.Tokens = append(c.class.Tokens, *close)

	return c.class, nil
}

// func (c *Compiler) CompileClass() (*Class, error) {

// }
