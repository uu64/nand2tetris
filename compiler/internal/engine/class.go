package engine

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/symtab"
	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

type Class struct {
	XMLName xml.Name `xml:"class"`
	Tokens  []tokenizer.Element
}

func (el Class) ElementType() tokenizer.ElementType {
	return tokenizer.ElClass
}

type ClassVarDec struct {
	XMLName xml.Name `xml:"classVarDec"`
	Tokens  []tokenizer.Element
}

func (el ClassVarDec) ElementType() tokenizer.ElementType {
	return tokenizer.ElClassVarDec
}

type SubroutineDec struct {
	XMLName xml.Name `xml:"subroutineDec"`
	Tokens  []tokenizer.Element
}

func (el SubroutineDec) ElementType() tokenizer.ElementType {
	return tokenizer.ElSubroutineDec
}

type ParameterList struct {
	XMLName xml.Name `xml:"parameterList"`
	Tokens  []tokenizer.Element
}

func (el ParameterList) ElementType() tokenizer.ElementType {
	return tokenizer.ElParameterList
}

type SubroutineBody struct {
	XMLName xml.Name `xml:"subroutineBody"`
	Tokens  []tokenizer.Element
}

func (el SubroutineBody) ElementType() tokenizer.ElementType {
	return tokenizer.ElSubroutineBody
}

type VarDec struct {
	XMLName xml.Name `xml:"varDec"`
	Tokens  []tokenizer.Element
}

func (el VarDec) ElementType() tokenizer.ElementType {
	return tokenizer.ElVarDec
}

func (c *Compiler) CompileClass() (*Class, error) {
	class := &Class{Tokens: []tokenizer.Element{}}

	if !c.tokenizer.HasMoreTokens() {
		return nil, fmt.Errorf("CompileClass: no tokens")
	}

	// 'class'
	kwd, err := c.consumeKeyword(tokenizer.KwdClass)
	if err != nil {
		return nil, fmt.Errorf("CompileClass: class should start with CLASS, got %v", c.tokenizer.Current)
	}
	class.Tokens = append(class.Tokens, kwd)

	// className
	className, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}
	class.Tokens = append(class.Tokens, className)

	c.ctx.ClassName = className.Val()
	className.Kind = symtab.SkNone.String()
	className.Category = symtab.SkClass.String()
	className.Index = -1
	className.IsDefined = true

	// '{'
	open, err := c.consumeSymbol(tokenizer.SymLeftCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileClass: symbol '{' is missing, got %v", c.tokenizer.Current)
	}
	class.Tokens = append(class.Tokens, open)

	// classVarDec* subroutineDec*
	for c.tokenizer.TokenType() == tokenizer.TkKeyword {
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()

		switch kwd.Val() {
		// classVarDec*
		case tokenizer.KwdStatic, tokenizer.KwdField:
			classVarDec, err := c.CompileClassVarDec()
			if err != nil {
				return nil, fmt.Errorf("CompileClass: %w", err)
			}
			class.Tokens = append(class.Tokens, classVarDec)
		// subroutineDec*
		case tokenizer.KwdConstructor, tokenizer.KwdFunction, tokenizer.KwdMethod:
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
	close, err := c.consumeSymbol(tokenizer.SymRightCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileClass: symbol '}' is missing, got %v", c.tokenizer.Current)
	}
	class.Tokens = append(class.Tokens, close)

	c.codewriter.Close()
	return class, nil
}

func (c *Compiler) CompileClassVarDec() (*ClassVarDec, error) {
	classVarDec := &ClassVarDec{Tokens: []tokenizer.Element{}}

	// ('static' | 'field')
	kwd, err := c.consumeKeyword(tokenizer.KwdStatic, tokenizer.KwdField)
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: classVarDec should start with STATIC or FIELD, got %v", c.tokenizer.Current)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, kwd)

	// type
	typ, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileClassVarDec: %w", err)
	}
	classVarDec.Tokens = append(classVarDec.Tokens, typ)

	// varName (',' varName)* ';'
	for {
		// varName
		varName, err := c.compileName()
		if err != nil {
			return nil, fmt.Errorf("CompileClassVarDec: %w", err)
		}
		classVarDec.Tokens = append(classVarDec.Tokens, varName)

		c.defineSymbol(varName, varName.Label, symtab.ElmToTyp(typ), symtab.KwdToKind(kwd))

		// check additional varName
		s, err := c.consumeSymbol(tokenizer.SymComma, tokenizer.SymSemiColon)
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: expected \",\" or \";\", got %v", s)
		}

		if s.Val() == tokenizer.SymComma {
			classVarDec.Tokens = append(classVarDec.Tokens, s)
		} else if s.Val() == tokenizer.SymSemiColon {
			classVarDec.Tokens = append(classVarDec.Tokens, s)
			break
		}
	}
	return classVarDec, nil
}

func (c *Compiler) CompileSubroutineDec() (*SubroutineDec, error) {
	subroutineDec := &SubroutineDec{Tokens: []tokenizer.Element{}}
	c.symtab.StartSubroutine()
	c.ctx.WhileIndex = 0
	c.ctx.IfIndex = 0

	// TODO: methodの場合だけ追加する
	// c.defineSymbol(&tokenizer.Identifier{Label: "this"}, "this", c.ctx.ClassName, symtab.SkArg)

	// ('constructor' | 'function' | 'method')
	kwd, err := c.consumeKeyword(tokenizer.KwdConstructor, tokenizer.KwdFunction, tokenizer.KwdMethod)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: classVarDec should start with CONSTRUCTOR or FUNCTION or METHOD, got %v", c.tokenizer.Current)
	} else {
		subroutineDec.Tokens = append(subroutineDec.Tokens, kwd)
	}
	c.ctx.SubroutineKwd = kwd

	// ('void' | type)
	isVoid := false
	if kwd, err := c.tokenizer.Keyword(); err == nil && kwd.Val() == tokenizer.KwdVoid {
		isVoid = true
		subroutineDec.Tokens = append(subroutineDec.Tokens, kwd)
		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	} else {
		typ, err := c.compileType()
		if err != nil {
			return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
		}
		subroutineDec.Tokens = append(subroutineDec.Tokens, typ)
	}

	// subroutineName
	subroutineName, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	c.ctx.SubroutineName = subroutineName.Val()
	subroutineDec.Tokens = append(subroutineDec.Tokens, subroutineName)
	subroutineName.Kind = symtab.SkNone.String()
	subroutineName.Category = symtab.SkSubroutine.String()
	subroutineName.Index = -1
	subroutineName.IsDefined = true

	// '('
	open, err := c.consumeSymbol(tokenizer.SymLeftParenthesis)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: symbol '(' is missing, got %v", c.tokenizer.Current)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, open)

	// parameterList
	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileParameterList()
	paramList, err := c.CompileParameterList()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, paramList)

	// ')'
	close, err := c.consumeSymbol(tokenizer.SymRightParenthesis)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: symbol ')' is missing, got %v", c.tokenizer.Current)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, close)

	// subroutineBody
	subroutineBody, err := c.CompileSubroutineBody()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineDec: %w", err)
	}
	subroutineDec.Tokens = append(subroutineDec.Tokens, subroutineBody)

	c.writeReturn(isVoid)
	return subroutineDec, nil
}

func (c *Compiler) CompileParameterList() (*ParameterList, error) {
	parameterList := &ParameterList{Tokens: []tokenizer.Element{}}

	// ((type varName) (',' type varName)*)
	consumeParameters := func() error {
		for {
			// type
			typ, err := c.compileType()
			if err != nil {
				return fmt.Errorf("CompileParameterList: %w", err)
			}
			parameterList.Tokens = append(parameterList.Tokens, typ)

			// varName
			varName, err := c.compileName()
			if err != nil {
				return fmt.Errorf("CompileParameterList: %w", err)
			}
			parameterList.Tokens = append(parameterList.Tokens, varName)
			c.defineSymbol(varName, varName.Label, symtab.ElmToTyp(typ), symtab.SkArg)

			// check additional parameter
			s, err := c.tokenizer.Symbol()
			if err != nil {
				return err
			}
			if s.Val() != tokenizer.SymComma {
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
	if c.tokenizer.TokenType() == tokenizer.TkKeyword {
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()
		if v := kwd.Val(); v == tokenizer.KwdInt || v == tokenizer.KwdChar || v == tokenizer.KwdBoolean {
			err := consumeParameters()
			if err != nil {
				return nil, err
			}
		}
	}
	if c.tokenizer.TokenType() == tokenizer.TkIdentifier {
		err := consumeParameters()
		if err != nil {
			return nil, err
		}
	}

	return parameterList, nil
}

func (c *Compiler) CompileSubroutineBody() (*SubroutineBody, error) {
	subroutineBody := &SubroutineBody{Tokens: []tokenizer.Element{}}

	// '{'
	open, err := c.consumeSymbol(tokenizer.SymLeftCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: symbol '{' is missing, got %v", c.tokenizer.Current)
	}
	subroutineBody.Tokens = append(subroutineBody.Tokens, open)

	// varDec*
	for {
		if kwd, err := c.tokenizer.Keyword(); !(err == nil && kwd.Val() == tokenizer.KwdVar) {
			break
		}
		varDec, err := c.CompileVarDec()
		if err != nil {
			return nil, fmt.Errorf("CompileSubroutineBody: %w", err)
		}
		subroutineBody.Tokens = append(subroutineBody.Tokens, varDec)
	}

	c.writeFunction(c.ctx.SubroutineKwd.Label, c.ctx.SubroutineName, c.symtab.VarCount(symtab.SkVar))

	// statements
	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()
	statements, err := c.CompileStatements()
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: %w", err)
	}
	subroutineBody.Tokens = append(subroutineBody.Tokens, statements)

	// '}'
	close, err := c.consumeSymbol(tokenizer.SymRightCurlyBracket)
	if err != nil {
		return nil, fmt.Errorf("CompileSubroutineBody: symbol '}' is missing, got %v", c.tokenizer.Current)
	}
	subroutineBody.Tokens = append(subroutineBody.Tokens, close)

	return subroutineBody, nil
}

func (c *Compiler) CompileVarDec() (*VarDec, error) {
	varDec := &VarDec{Tokens: []tokenizer.Element{}}

	// 'var'
	kwd, err := c.consumeKeyword(tokenizer.KwdVar)
	if err != nil {
		return nil, fmt.Errorf("CompileVarDec: varDec should start with VAR, got %v", c.tokenizer.Current)
	}
	varDec.Tokens = append(varDec.Tokens, kwd)

	// type
	typ, err := c.compileType()
	if err != nil {
		return nil, fmt.Errorf("CompileVarDec: %w", err)
	}
	varDec.Tokens = append(varDec.Tokens, typ)

	// varName (',' varName)*
	for {
		// varName
		varName, err := c.compileName()
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: %w", err)
		}
		varDec.Tokens = append(varDec.Tokens, varName)
		c.defineSymbol(varName, varName.Label, symtab.ElmToTyp(typ), symtab.SkVar)

		// check additional varName
		s, err := c.consumeSymbol(tokenizer.SymComma, tokenizer.SymSemiColon)
		if err != nil {
			return nil, fmt.Errorf("CompileVarDec: expected \",\" or \";\", got %v", s)
		}

		if s.Val() == tokenizer.SymComma {
			varDec.Tokens = append(varDec.Tokens, s)
		} else if s.Val() == tokenizer.SymSemiColon {
			varDec.Tokens = append(varDec.Tokens, s)
			break
		}
	}
	return varDec, nil
}

func (c *Compiler) compileType() (tokenizer.Element, error) {
	tkType := c.tokenizer.Current.TokenType()
	switch tkType {
	case tokenizer.TkKeyword:
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()
		if kwd.Val() != tokenizer.KwdInt && kwd.Val() != tokenizer.KwdChar && kwd.Val() != tokenizer.KwdBoolean {
			return nil, fmt.Errorf("compileType: invalid type %s", kwd.Label)
		}
		return kwd, c.tokenizer.Advance()
	case tokenizer.TkIdentifier:
		// ignore the error because it is already checked that the token type is IDENTIFIER
		return c.compileName()
	default:
		return nil, fmt.Errorf("compileType: type should start with KEYWORD or IDENTIFIER, got %d", tkType)
	}
}

func (c *Compiler) compileName() (*tokenizer.Identifier, error) {
	id, err := c.tokenizer.Identifier()

	name := id.Label
	kind := c.symtab.KindOf(name)

	if kind != symtab.SkNone {
		id.Kind = kind.String()
		id.Category = kind.String()
		id.Index = c.symtab.IndexOf(name)
		id.IsDefined = false
	} else {
		if name == c.ctx.ClassName {
			id.Kind = kind.String()
			id.Category = symtab.SkClass.String()
			id.Index = c.symtab.IndexOf(name)
			id.IsDefined = false
		} else if name == c.ctx.SubroutineName {
			id.Kind = kind.String()
			id.Category = symtab.SkSubroutine.String()
			id.Index = c.symtab.IndexOf(name)
			id.IsDefined = false
		} else {
			id.Kind = kind.String()
			id.Category = kind.String()
			id.Index = c.symtab.IndexOf(name)
			id.IsDefined = true
		}
	}

	// fmt.Printf("id name: %s\n", name)
	// fmt.Printf("class name: %s\n", c.symtab.ClassName)
	// fmt.Printf("subroutine name: %s\n", c.symtab.SubroutineName)
	// fmt.Printf("class: %v\n", c.symtab.ClassTable())
	// fmt.Printf("subroutine: %v\n", c.symtab.SubroutineTable())
	// fmt.Println()

	if err != nil {
		return nil, fmt.Errorf("compileName: %w", err)
	}

	return id, c.tokenizer.Advance()
}

func (c *Compiler) defineSymbol(id *tokenizer.Identifier, symName, typ string, kind symtab.SymbolKind) {
	c.symtab.Define(symName, typ, kind)
	id.Kind = c.symtab.KindOf(symName).String()
	id.Category = c.symtab.KindOf(symName).String()
	id.Index = c.symtab.IndexOf(symName)
	id.IsDefined = true
}
