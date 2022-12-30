package compile

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Class struct {
	XMLName xml.Name `xml:"class"`
	Tokens  []token.Element
}

func (el Class) ElementType() token.ElementType {
	return token.ElClass
}

type ClassVarDec struct {
	XMLName xml.Name `xml:"classVarDec"`
	Tokens  []token.Element
}

func (el ClassVarDec) ElementType() token.ElementType {
	return token.ElClassVarDec
}

type SubroutineDec struct {
	XMLName xml.Name `xml:"subroutineDec"`
	Tokens  []token.Element
}

func (el SubroutineDec) ElementType() token.ElementType {
	return token.ElSubroutineDec
}

type ParameterList struct {
	XMLName xml.Name `xml:"parameterList"`
	Tokens  []token.Element
}

func (el ParameterList) ElementType() token.ElementType {
	return token.ElParameterList
}

type SubroutineBody struct {
	XMLName xml.Name `xml:"subroutineBody"`
	Tokens  []token.Element
}

func (el SubroutineBody) ElementType() token.ElementType {
	return token.ElSubroutineBody
}

type VarDec struct {
	XMLName xml.Name `xml:"varDec"`
	Tokens  []token.Element
}

func (el VarDec) ElementType() token.ElementType {
	return token.ElVarDec
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
	className, err := c.compileClassName()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	class.Tokens = append(class.Tokens, className)

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
			subroutineDec, err := c.CompileSubroutineDec()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			class.Tokens = append(class.Tokens, subroutineDec)
		default:
			return nil, fmt.Errorf("CompileClass: invalid token %v", kwd)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
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
		classVarDec.Tokens = append(classVarDec.Tokens, varName)

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
			break
		} else {
			return nil, fmt.Errorf("CompileClassVarDec: expected \",\" or \";\", got %v", s)
		}
	}
	return classVarDec, nil
}

func (c *Compiler) CompileSubroutineDec() (*SubroutineDec, error) {
	subroutineDec := &SubroutineDec{Tokens: []token.Element{}}

	// constructor | function | method で始まるかチェック
	if kwd, err := c.tokenizer.Keyword(); err != nil || (kwd.Val() != token.KwdConstructor && kwd.Val() != token.KwdFunction && kwd.Val() != token.KwdMethod) {
		return nil, fmt.Errorf("CompileSubroutineDec: classVarDec should start with CONSTRUCTOR or FUNCTION or METHOD, got %v", c.tokenizer.Current)
	} else {
		subroutineDec.Tokens = append(subroutineDec.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// void | type
	kwd, err := c.tokenizer.Keyword()
	if err == nil && kwd.Val() == token.KwdVoid {
		subroutineDec.Tokens = append(subroutineDec.Tokens, *kwd)
	} else {
		types, err := c.compileType()
		if err != nil {
			return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
		}
		subroutineDec.Tokens = append(subroutineDec.Tokens, types...)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// subroutine name
	subroutineName, err := c.compileSubroutineName()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, subroutineName)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ( があるかチェック
	if open, err := c.tokenizer.Symbol(); err != nil || open.Val() != rune('(') {
		return nil, fmt.Errorf("CompileSubroutineDec: symbol '(' is missing, got %v", c.tokenizer.Current)
	} else {
		subroutineDec.Tokens = append(subroutineDec.Tokens, *open)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// parameter list
	paramList, err := c.CompileParameterList()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, paramList)

	// NOTE: CompileParameterListで先読みされるためAdvanceの呼び出しは不要

	// ) があるかチェック
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune(')') {
		return nil, fmt.Errorf("CompileSubroutineDec: symbol ')' is missing, got %v", c.tokenizer.Current)
	} else {
		subroutineDec.Tokens = append(subroutineDec.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// subroutine body
	subroutineBody, err := c.CompileSubroutineBody()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, subroutineBody)

	return subroutineDec, nil
}

func (c *Compiler) CompileParameterList() (*ParameterList, error) {
	parameterList := &ParameterList{Tokens: []token.Element{}}

	// 現在のtokenがidentifierではない場合、空のparameterListを返す
	if c.tokenizer.TokenType() != token.TkKeyword {
		return parameterList, nil
	}

	for {
		// type
		types, err := c.compileType()
		if err != nil {
			return nil, fmt.Errorf("CompileParameterList: %w", err)
		}
		parameterList.Tokens = append(parameterList.Tokens, types...)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// varName
		varName, err := c.compileVarName()
		if err != nil {
			return nil, fmt.Errorf("CompileParameterList: %w", err)
		}
		parameterList.Tokens = append(parameterList.Tokens, varName)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// check additional parameter
		if s, err := c.tokenizer.Symbol(); err != nil || s.Val() != rune(',') {
			break
		} else {
			parameterList.Tokens = append(parameterList.Tokens, s)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	}

	return parameterList, nil
}

func (c *Compiler) CompileSubroutineBody() (*SubroutineBody, error) {
	subroutineBody := &SubroutineBody{Tokens: []token.Element{}}

	// { があるかチェック
	if open, err := c.tokenizer.Symbol(); err != nil || open.Val() != rune('{') {
		return nil, fmt.Errorf("CompileSubroutineBody: symbol '{' is missing, got %v", c.tokenizer.Current)
	} else {
		subroutineBody.Tokens = append(subroutineBody.Tokens, *open)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	for c.tokenizer.TokenType() == token.TkKeyword {
		// TkKeywordであると確定しているためerrorは起こり得ない
		kwd, _ := c.tokenizer.Keyword()

		switch kwd.Val() {
		case token.KwdVar:
			classVarDec, err := c.CompileClassVarDec()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			subroutineBody.Tokens = append(subroutineBody.Tokens, classVarDec)
		case token.KwdLet, token.KwdIf, token.KwdWhile, token.KwdDo, token.KwdReturn:
			subroutineDec, err := c.CompileStatements()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			subroutineBody.Tokens = append(subroutineBody.Tokens, subroutineDec)
		default:
			return nil, fmt.Errorf("CompileClass: invalid token %v", kwd)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	}

	// } があるかチェック
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('}') {
		return nil, fmt.Errorf("CompileSubroutineBody: symbol '}' is missing, got %v", c.tokenizer.Current)
	} else {
		subroutineBody.Tokens = append(subroutineBody.Tokens, *close)
	}

	return subroutineBody, nil
}

func (c *Compiler) CompileVarDec() (*VarDec, error) {
	varDec := &VarDec{Tokens: []token.Element{}}

	// var で始まるかチェック
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdVar {
		return nil, fmt.Errorf("CompileVarDec: varDec should start with VAR, got %v", c.tokenizer.Current)
	} else {
		varDec.Tokens = append(varDec.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// type
	types, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileVarDec: %w", err)
	}
	varDec.Tokens = append(varDec.Tokens, types...)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	for {
		// varName
		varName, err := c.compileVarName()
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: %w", err)
		}
		varDec.Tokens = append(varDec.Tokens, varName)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// check additional varName
		s, err := c.tokenizer.Symbol()
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: expected \",\" or \";\", got %v", s)
		}

		if s.Val() == rune(',') {
			varDec.Tokens = append(varDec.Tokens, *s)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, err
			}
		} else if s.Val() == rune(';') {
			varDec.Tokens = append(varDec.Tokens, *s)
			break
		} else {
			return nil, fmt.Errorf("CompileVarDec: expected \",\" or \";\", got %v", s)
		}
	}
	return varDec, nil
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

func (c *Compiler) compileClassName() (*token.Identifier, error) {
	return c.compileName()
}

func (c *Compiler) compileSubroutineName() (*token.Identifier, error) {
	return c.compileName()
}

func (c *Compiler) compileVarName() (*token.Identifier, error) {
	return c.compileName()
}

func (c *Compiler) compileName() (*token.Identifier, error) {
	id, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("compileVarName: %w", err)
	}
	return id, nil
}
