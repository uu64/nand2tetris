package compile

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
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

func (c *Compiler) CompileClass() (*Class, error) {
	class := &Class{Tokens: []token.Element{}}

	if !c.tokenizer.HasMoreTokens() {
		return nil, fmt.Errorf("CompileClass: no tokens")
	}

	// class で始まるかチェック
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdClass {
		return nil, fmt.Errorf("CompileClass: class should start with CLASS, got %v", c.tokenizer.Current)
	} else {
		class.Tokens = append(class.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}

	// className があるかチェック
	if className, err := c.tokenizer.Identifier(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	} else {
		class.Tokens = append(class.Tokens, *className)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}

	// { があるかチェック
	if open, err := c.tokenizer.Symbol(); err != nil || open.Val() != rune('{') {
		return nil, fmt.Errorf("CompileClass: symbol '{' is missing, got %v", c.tokenizer.Current)
	} else {
		class.Tokens = append(class.Tokens, *open)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	for c.tokenizer.TokenType() == token.TkKeyword {
		// TkKeywordであると確定しているためerrorは起こり得ない
		kwd, _ := c.tokenizer.Keyword()
		switch kwd.Val() {
		case token.KwdStatic, token.KwdField:
			classVarDec, err := c.CompileClassVarDec()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			class.Tokens = append(class.Tokens, classVarDec)
		case token.KwdConstructor, token.KwdFunction, token.KwdMethod:
			// continue
		default:
			return nil, fmt.Errorf("CompileClass: invalid token %v", kwd)
		}
	}

	// } があるかチェック
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('}') {
		return nil, fmt.Errorf("CompileClass: symbol '}' is missing, got %v", c.tokenizer.Current)
	} else {
		class.Tokens = append(class.Tokens, *close)
	}

	return class, nil
}

func (c *Compiler) CompileClassVarDec() (*ClassVarDec, error) {
	classVarDec := &ClassVarDec{Tokens: []token.Element{}}

	// static | field で始まるかチェック
	if kwd, err := c.tokenizer.Keyword(); err != nil || (kwd.Val() != token.KwdStatic && kwd.Val() != token.KwdField) {
		return nil, fmt.Errorf("CompileClassVarDec: classVarDec should start with STATIC or FIELD, got %v", c.tokenizer.Current)
	} else {
		classVarDec.Tokens = append(classVarDec.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// type
	types, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, types...)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	for {
		// varName
		varName, err := c.compileVarName()
		if err != nil {
			return nil, fmt.Errorf("CompileClassVarDec: %w", err)
		}
		classVarDec.Tokens = append(classVarDec.Tokens, varName...)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// check additional varName
		s, err := c.tokenizer.Symbol()
		if err != nil {
			return nil, fmt.Errorf("CompileClassVarDec: expected \",\" or \";\", got %v", s)
		}

		if s.Val() == rune(',') {
			classVarDec.Tokens = append(classVarDec.Tokens, *s)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, err
			}
		} else if s.Val() == rune(';') {
			classVarDec.Tokens = append(classVarDec.Tokens, *s)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, err
			}
			break
		} else {
			return nil, fmt.Errorf("CompileClassVarDec: expected \",\" or \";\", got %v", s)
		}
	}
	return classVarDec, nil
}

func (c *Compiler) compileType() ([]token.Element, error) {
	tokens := []token.Element{}

	tkType := c.tokenizer.Current.TokenType()
	switch tkType {
	case token.TkKeyword:
		// TkKeywordであると確定しているためerrorは起こり得ない
		kwd, _ := c.tokenizer.Keyword()
		if kwd.Val() != token.KwdInt && kwd.Val() != token.KwdChar && kwd.Val() != token.KwdBoolean {
			return nil, fmt.Errorf("compileType: invalid type %s", kwd.Label)
		}
		tokens = append(tokens, *kwd)
	case token.TkIdentifier:
		// TkIdentifierであると確定しているためerrorは起こり得ない
		id, _ := c.tokenizer.Identifier()
		tokens = append(tokens, *id)
	default:
		return nil, fmt.Errorf("compileType: type should start with keyword or identifier, got %d", tkType)
	}

	return tokens, nil
}

func (c *Compiler) compileVarName() ([]token.Element, error) {
	tokens := []token.Element{}

	id, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("compileVarName: %w", err)
	}
	tokens = append(tokens, *id)
	return tokens, nil
}
