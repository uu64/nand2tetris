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
	class := &Class{Tokens: []token.Element{}}

	if !c.tokenizer.HasMoreTokens() {
		return nil, fmt.Errorf("CompileClass: no tokens")
	}

	// class で始まるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	kwd, err := c.tokenizer.Keyword()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	if kwd.Val() != token.KwdClass {
		return nil, fmt.Errorf("CompileClass: class should start with CLASS, got %s", kwd.Label)
	}
	class.Tokens = append(class.Tokens, *kwd)

	// className があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	className, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	class.Tokens = append(class.Tokens, *className)

	// { があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}
	open, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	if open.Val() != rune('{') {
		return nil, fmt.Errorf("CompileClass: symbol '{' is missing")
	}
	class.Tokens = append(class.Tokens, *open)

	// classVarDec
	classVarDec, err := c.CompileClassVarDec()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	class.Tokens = append(class.Tokens, classVarDec)

	// c.CompileSubroutineVarDec()

	// } があるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	close, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	if close.Val() != rune('}') {
		return nil, fmt.Errorf("CompileClass: symbol '}' is missing")
	}
	class.Tokens = append(class.Tokens, *close)

	return class, nil
}

func (c *Compiler) CompileClassVarDec() (*ClassVarDec, error) {
	classVarDec := &ClassVarDec{Tokens: []token.Element{}}

	// static | field で始まるかチェック
	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	kwd, err := c.tokenizer.Keyword()
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	if v := kwd.Val(); v != token.KwdStatic && v != token.KwdField {
		return nil, fmt.Errorf("CompileClassVarDec: classVarDec should start with static or field, got %s", kwd.Label)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, *kwd)

	// type
	types, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, types...)

	// varName
	varName, err := c.compileVarName()
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, varName...)

	for {
		if err := c.tokenizer.Advance(); err != nil {
			return nil, fmt.Errorf("CompileClassVarDec: %w", err)
		}

		symbol, err := c.tokenizer.Symbol()
		if err != nil {
			return nil, fmt.Errorf("CompileClassVarDec: symbol is expected: %w", err)
		}

		if symbol.Val() == rune(';') {
			classVarDec.Tokens = append(classVarDec.Tokens, symbol)
			return classVarDec, nil
		}

		if symbol.Val() == rune(',') {
			classVarDec.Tokens = append(classVarDec.Tokens, symbol)
			varName, err := c.compileVarName()
			if err != nil {
				return nil, fmt.Errorf("CompileClassVarDec: %w", err)
			}
			classVarDec.Tokens = append(classVarDec.Tokens, varName...)
		} else {
			return nil, fmt.Errorf("CompileClassVarDec: unexpected symbol %s", symbol.Label)
		}
	}
}

func (c *Compiler) compileType() ([]token.Element, error) {
	tokens := []token.Element{}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("compileType: %w", err)
	}
	switch c.tokenizer.TkType {
	case token.TkKeyword:
		kwd, err := c.tokenizer.Keyword()
		if err != nil {
			return nil, fmt.Errorf("compileType: %w", err)
		}
		if v := kwd.Val(); v != token.KwdInt && v != token.KwdChar && v != token.KwdBoolean {
			return nil, fmt.Errorf("compileType: invalid type %s", kwd.Label)
		}
		tokens = append(tokens, *kwd)
	case token.TkIdentifier:
		id, err := c.tokenizer.Identifier()
		if err != nil {
			return nil, fmt.Errorf("compileType: %w", err)
		}
		tokens = append(tokens, *id)
	default:
		return nil, fmt.Errorf("compileType: type should start with keyword or identifier, got %d", c.tokenizer.TkType)
	}

	return tokens, nil
}

func (c *Compiler) compileVarName() ([]token.Element, error) {
	tokens := []token.Element{}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("compileVarName: %w", err)
	}

	id, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("compileVarName: %w", err)
	}
	tokens = append(tokens, *id)
	return tokens, nil
}
