package engine

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/symtab"
	token "github.com/uu64/nand2tetris/compiler/internal/tokenizer"
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

	// 'class'
	kwd, err := c.consumeKeyword(token.KwdClass)
	if err != nil {
		return nil, fmt.Errorf("CompileClass: class should start with CLASS, got %v", c.tokenizer.Current)
	}
	class.Tokens = append(class.Tokens, *kwd)

	// className
	className, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	class.Tokens = append(class.Tokens, className)
	c.symtab.ClassName = className.Val()

	// '{'
	open, err := c.consumeSymbol(token.SymLeftCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileClass: symbol '{' is missing, got %v", c.tokenizer.Current)
	}
	class.Tokens = append(class.Tokens, *open)

	// classVarDec* subroutineDec*
	for c.tokenizer.TokenType() == token.TkKeyword {
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()

		switch kwd.Val() {
		// classVarDec*
		case token.KwdStatic, token.KwdField:
			classVarDec, err := c.CompileClassVarDec()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			class.Tokens = append(class.Tokens, classVarDec)
		// subroutineDec*
		case token.KwdConstructor, token.KwdFunction, token.KwdMethod:
			subroutineDec, err := c.CompileSubroutineDec()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			class.Tokens = append(class.Tokens, subroutineDec)
		default:
			return nil, fmt.Errorf("CompileClass: invalid token %v", kwd)
		}
	}

	// '}'
	close, err := c.consumeSymbol(token.SymRightCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileClass: symbol '}' is missing, got %v", c.tokenizer.Current)
	}
	class.Tokens = append(class.Tokens, *close)

	// fmt.Println(c.symtab.ClassTable())

	return class, nil
}

func (c *Compiler) CompileClassVarDec() (*ClassVarDec, error) {
	classVarDec := &ClassVarDec{Tokens: []token.Element{}}

	// ('static' | 'field')
	kwd, err := c.consumeKeyword(token.KwdStatic, token.KwdField)
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: classVarDec should start with STATIC or FIELD, got %v", c.tokenizer.Current)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, *kwd)

	// type
	typ, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, typ)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// varName (',' varName)* ';'
	for {
		// varName
		varName, err := c.compileName()
		if err != nil {
			return nil, fmt.Errorf("CompileClassVarDec: %w", err)
		}
		classVarDec.Tokens = append(classVarDec.Tokens, varName)

		c.symtab.Define(varName.Label, symtab.ElmToTyp(typ), symtab.KwdToKind(kwd))

		// check additional varName
		s, err := c.consumeSymbol(token.SymComma, token.SymSemiColon)
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: expected \",\" or \";\", got %v", s)
		}

		if s.Val() == token.SymComma {
			classVarDec.Tokens = append(classVarDec.Tokens, *s)
		} else if s.Val() == token.SymSemiColon {
			classVarDec.Tokens = append(classVarDec.Tokens, *s)
			break
		}
	}
	return classVarDec, nil
}

func (c *Compiler) CompileSubroutineDec() (*SubroutineDec, error) {
	subroutineDec := &SubroutineDec{Tokens: []token.Element{}}
	c.symtab.StartSubroutine()
	c.symtab.Define("this", c.symtab.ClassName, symtab.SkArg)

	// ('constructor' | 'function' | 'method')
	if kwd, err := c.consumeKeyword(token.KwdConstructor, token.KwdFunction, token.KwdMethod); err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: classVarDec should start with CONSTRUCTOR or FUNCTION or METHOD, got %v", c.tokenizer.Current)
	} else {
		subroutineDec.Tokens = append(subroutineDec.Tokens, *kwd)
	}

	// ('void' | type)
	if kwd, err := c.tokenizer.Keyword(); err == nil && kwd.Val() == token.KwdVoid {
		subroutineDec.Tokens = append(subroutineDec.Tokens, *kwd)
	} else {
		typ, err := c.compileType()
		if err != nil {
			return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
		}
		subroutineDec.Tokens = append(subroutineDec.Tokens, typ)
	}
	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// subroutineName
	subroutineName, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, subroutineName)

	// '('
	open, err := c.consumeSymbol(token.SymLeftParenthesis)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: symbol '(' is missing, got %v", c.tokenizer.Current)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, *open)

	// parameterList
	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileParameterList()
	paramList, err := c.CompileParameterList()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, paramList)

	// ')'
	close, err := c.consumeSymbol(token.SymRightParenthesis)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: symbol ')' is missing, got %v", c.tokenizer.Current)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, *close)

	// subroutineBody
	subroutineBody, err := c.CompileSubroutineBody()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, subroutineBody)

	// fmt.Println(c.symtab.SubroutineTable())
	return subroutineDec, nil
}

func (c *Compiler) CompileParameterList() (*ParameterList, error) {
	parameterList := &ParameterList{Tokens: []token.Element{}}

	// ((type varName) (',' type varName)*)
	consumeParameters := func() error {
		for {
			// type
			typ, err := c.compileType()
			if err != nil {
				return fmt.Errorf("CompileParameterList: %w", err)
			}
			parameterList.Tokens = append(parameterList.Tokens, typ)

			if err := c.tokenizer.Advance(); err != nil {
				return err
			}

			// varName
			varName, err := c.compileName()
			if err != nil {
				return fmt.Errorf("CompileParameterList: %w", err)
			}
			parameterList.Tokens = append(parameterList.Tokens, varName)
			c.symtab.Define(varName.Label, symtab.ElmToTyp(typ), symtab.SkArg)

			// check additional parameter
			s, err := c.tokenizer.Symbol()
			if err != nil {
				return err
			}
			if s.Val() != token.SymComma {
				break
			}
			parameterList.Tokens = append(parameterList.Tokens, s)
			if err := c.tokenizer.Advance(); err != nil {
				return err
			}
		}
		return nil
	}

	// Return empty ParameterList when current token is not type
	if c.tokenizer.TokenType() == token.TkKeyword {
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()
		if v := kwd.Val(); v == token.KwdInt || v == token.KwdChar || v == token.KwdBoolean {
			err := consumeParameters()
			if err != nil {
				return nil, err
			}
		}
	}
	if c.tokenizer.TokenType() == token.TkIdentifier {
		err := consumeParameters()
		if err != nil {
			return nil, err
		}
	}

	return parameterList, nil
}

func (c *Compiler) CompileSubroutineBody() (*SubroutineBody, error) {
	subroutineBody := &SubroutineBody{Tokens: []token.Element{}}

	// '{'
	open, err := c.consumeSymbol(token.SymLeftCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: symbol '{' is missing, got %v", c.tokenizer.Current)
	}
	subroutineBody.Tokens = append(subroutineBody.Tokens, *open)

	// varDec*
	for {
		if kwd, err := c.tokenizer.Keyword(); !(err == nil && kwd.Val() == token.KwdVar) {
			break
		}
		varDec, err := c.CompileVarDec()
		if err != nil {
			return nil, fmt.Errorf("CompileSubroutineBody: %w", err)
		}
		subroutineBody.Tokens = append(subroutineBody.Tokens, varDec)
	}

	// statements
	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()
	statements, err := c.CompileStatements()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: %w", err)
	}
	subroutineBody.Tokens = append(subroutineBody.Tokens, statements)

	// '}'
	close, err := c.consumeSymbol(token.SymRightCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: symbol '}' is missing, got %v", c.tokenizer.Current)
	}
	subroutineBody.Tokens = append(subroutineBody.Tokens, *close)

	return subroutineBody, nil
}

func (c *Compiler) CompileVarDec() (*VarDec, error) {
	varDec := &VarDec{Tokens: []token.Element{}}

	// 'var'
	kwd, err := c.consumeKeyword(token.KwdVar)
	if err != nil {
		return nil, fmt.Errorf("CompileVarDec: varDec should start with VAR, got %v", c.tokenizer.Current)
	}
	varDec.Tokens = append(varDec.Tokens, *kwd)

	// type
	typ, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileVarDec: %w", err)
	}
	varDec.Tokens = append(varDec.Tokens, typ)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// varName (',' varName)*
	for {
		// varName
		varName, err := c.compileName()
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: %w", err)
		}
		varDec.Tokens = append(varDec.Tokens, varName)
		c.symtab.Define(varName.Label, symtab.ElmToTyp(typ), symtab.SkVar)

		// check additional varName
		s, err := c.consumeSymbol(token.SymComma, token.SymSemiColon)
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: expected \",\" or \";\", got %v", s)
		}

		if s.Val() == token.SymComma {
			varDec.Tokens = append(varDec.Tokens, *s)
		} else if s.Val() == token.SymSemiColon {
			varDec.Tokens = append(varDec.Tokens, *s)
			break
		}
	}
	return varDec, nil
}

func (c *Compiler) compileType() (token.Element, error) {
	tkType := c.tokenizer.Current.TokenType()
	switch tkType {
	case token.TkKeyword:
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()
		if kwd.Val() != token.KwdInt && kwd.Val() != token.KwdChar && kwd.Val() != token.KwdBoolean {
			return nil, fmt.Errorf("compileType: invalid type %s", kwd.Label)
		}
		return kwd, nil
	case token.TkIdentifier:
		// ignore the error because it is already checked that the token type is IDENTIFIER
		id, _ := c.tokenizer.Identifier()
		return id, nil
	default:
		return nil, fmt.Errorf("compileType: type should start with KEYWORD or IDENTIFIER, got %d", tkType)
	}
}

func (c *Compiler) compileName() (*token.Identifier, error) {
	id, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("compileName: %w", err)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	return id, nil
}
