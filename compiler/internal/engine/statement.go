package engine

import (
	"encoding/xml"
	"fmt"

	token "github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

type Statements struct {
	XMLName xml.Name `xml:"statements"`
	Tokens  []token.Element
}

func (el Statements) ElementType() token.ElementType {
	return token.ElStatements
}

type LetStatement struct {
	XMLName xml.Name `xml:"letStatement"`
	Tokens  []token.Element
}

func (el LetStatement) ElementType() token.ElementType {
	return token.ElLetStatement
}

type IfStatement struct {
	XMLName xml.Name `xml:"ifStatement"`
	Tokens  []token.Element
}

func (el IfStatement) ElementType() token.ElementType {
	return token.ElIfStatement
}

type WhileStatement struct {
	XMLName xml.Name `xml:"whileStatement"`
	Tokens  []token.Element
}

func (el WhileStatement) ElementType() token.ElementType {
	return token.ElWhileStatement
}

type DoStatement struct {
	XMLName xml.Name `xml:"doStatement"`
	Tokens  []token.Element
}

func (el DoStatement) ElementType() token.ElementType {
	return token.ElDoStatement
}

type ReturnStatement struct {
	XMLName xml.Name `xml:"returnStatement"`
	Tokens  []token.Element
}

func (el ReturnStatement) ElementType() token.ElementType {
	return token.ElReturnStatement
}

func (c *Compiler) CompileStatements() (*Statements, error) {
	statements := Statements{Tokens: []token.Element{}}

	for c.tokenizer.TokenType() == token.TkKeyword {
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()

		switch kwd.Val() {
		case token.KwdLet:
			statement, err := c.compileLetStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case token.KwdIf:
			statement, err := c.compileIfStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case token.KwdWhile:
			statement, err := c.compileWhileStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case token.KwdDo:
			statement, err := c.compileDoStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case token.KwdReturn:
			statement, err := c.compileReturnStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		default:
			return nil, fmt.Errorf("CompileStatements: invalid keyword %s", kwd)
		}
	}

	return &statements, nil
}

func (c *Compiler) compileLetStatement() (*LetStatement, error) {
	statement := LetStatement{Tokens: []token.Element{}}

	// 'let'
	kwd, err := c.consumeKeyword(token.KwdLet)
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: letStatement should start with LET, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *kwd)

	// varName
	varName, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, varName)

	// ('[' expression ']')?
	if open, err := c.consumeSymbol(token.SymLeftSquareBracket); err == nil {
		// '['
		statement.Tokens = append(statement.Tokens, *open)

		// expression
		exp, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("compileLetStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, exp)

		// ']'
		close, err := c.consumeSymbol(token.SymRightSquareBracket)
		if err != nil {
			return nil, fmt.Errorf("compileLetStatement: symbol ']' is missing, got %v", c.tokenizer.Current)
		}
		statement.Tokens = append(statement.Tokens, *close)
	}

	// '='
	eq, err := c.consumeSymbol(token.SymEqual)
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: symbol '=' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *eq)

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	// ';'
	end, err := c.consumeSymbol(token.SymSemiColon)
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *end)

	return &statement, nil
}

func (c *Compiler) compileIfStatement() (*IfStatement, error) {
	statement := IfStatement{Tokens: []token.Element{}}

	// 'if'
	kwd, err := c.consumeKeyword(token.KwdIf)
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: ifStatement should start with IF, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *kwd)

	// '('
	open, err := c.consumeSymbol(token.SymLeftParenthesis)
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: symbol '(' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *open)

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	// ')'
	close, err := c.consumeSymbol(token.SymRightParenthesis)
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: symbol ')' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *close)

	consumeStatements := func() error {
		// '{'
		open, err := c.consumeSymbol(token.SymLeftCurlyBracket)
		if err != nil {
			return fmt.Errorf("compileIfStatement: symbol '{' is missing, got %v", c.tokenizer.Current)
		}
		statement.Tokens = append(statement.Tokens, *open)

		// statements
		// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()
		statements, err := c.CompileStatements()
		if err != nil {
			return fmt.Errorf("compileIfStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, statements)

		// '}'
		close, err := c.consumeSymbol(token.SymRightCurlyBracket)
		if err != nil {
			return fmt.Errorf("compileIfStatement: symbol '}' is missing, got %v", c.tokenizer.Current)
		}
		statement.Tokens = append(statement.Tokens, *close)
		return nil
	}

	if err := consumeStatements(); err != nil {
		return nil, err
	}

	// 'else'
	if kwd, err := c.consumeKeyword(token.KwdElse); err == nil {
		statement.Tokens = append(statement.Tokens, *kwd)

		if err := consumeStatements(); err != nil {
			return nil, err
		}
	}

	return &statement, nil
}

func (c *Compiler) compileWhileStatement() (*WhileStatement, error) {
	statement := WhileStatement{Tokens: []token.Element{}}

	// 'while'
	kwd, err := c.consumeKeyword(token.KwdWhile)
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: whileStatement should start with WHILE, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, *kwd)

	// '('
	if open, err := c.consumeSymbol(token.SymLeftParenthesis); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol '(' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *open)
	}

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	// ')'
	if close, err := c.consumeSymbol(token.SymRightParenthesis); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol ')' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	// '{'
	if open, err := c.consumeSymbol(token.SymLeftCurlyBracket); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol '{' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *open)
	}

	// statements
	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()
	statements, err := c.CompileStatements()
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, statements)

	// '}'
	if close, err := c.consumeSymbol(token.SymRightCurlyBracket); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol '}' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	return &statement, nil
}

func (c *Compiler) compileDoStatement() (*DoStatement, error) {
	statement := DoStatement{Tokens: []token.Element{}}

	// 'do'
	if kwd, err := c.consumeKeyword(token.KwdDo); err != nil {
		return nil, fmt.Errorf("compileDoStatement: doStatement should start with DO, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	// subroutineCall
	call, err := c.CompileSubroutineCall()
	if err != nil {
		return nil, fmt.Errorf("compileDoStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, call...)

	// ';'
	if end, err := c.consumeSymbol(token.SymSemiColon); err != nil {
		return nil, fmt.Errorf("compileDoStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *end)
	}

	return &statement, nil
}

func (c *Compiler) compileReturnStatement() (*ReturnStatement, error) {
	statement := ReturnStatement{Tokens: []token.Element{}}

	// 'return'
	if kwd, err := c.consumeKeyword(token.KwdReturn); err != nil {
		return nil, fmt.Errorf("compileReturnStatement: ReturnStatement should start with RETURN, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	// expression?
	if end, err := c.tokenizer.Symbol(); !(err == nil && end.Val() == token.SymSemiColon) {
		exp, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("compileReturnStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, exp)
	}

	// ';'
	if end, err := c.consumeSymbol(token.SymSemiColon); err != nil {
		return nil, fmt.Errorf("compileReturnStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *end)
	}

	return &statement, nil
}
