package compile

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
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
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdLet {
		return nil, fmt.Errorf("compileLetStatement: letStatement should start with LET, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// varName
	varName, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, varName)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ('[' expression ']')?
	if open, err := c.tokenizer.Symbol(); err == nil && open.Val() == rune('[') {
		// '['
		statement.Tokens = append(statement.Tokens, *open)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// expression
		exp, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("compileLetStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, exp)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// ']'
		if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune(']') {
			return nil, fmt.Errorf("compileLetStatement: symbol ']' is missing, got %v", c.tokenizer.Current)
		} else {
			statement.Tokens = append(statement.Tokens, *close)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	}

	// '='
	if eq, err := c.tokenizer.Symbol(); err != nil || eq.Val() != rune('=') {
		return nil, fmt.Errorf("compileLetStatement: symbol '=' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *eq)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ';'
	if end, err := c.tokenizer.Symbol(); err != nil || end.Val() != rune(';') {
		return nil, fmt.Errorf("compileLetStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *end)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	return &statement, nil
}

func (c *Compiler) compileIfStatement() (*IfStatement, error) {
	statement := IfStatement{Tokens: []token.Element{}}

	// 'if'
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdIf {
		return nil, fmt.Errorf("compileIfStatement: ifStatement should start with IF, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// '('
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('(') {
		return nil, fmt.Errorf("compileIfStatement: symbol '(' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ')'
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune(')') {
		return nil, fmt.Errorf("compileIfStatement: symbol ')' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	consumeStatements := func() error {
		// '{'
		if open, err := c.tokenizer.Symbol(); err != nil || open.Val() != rune('{') {
			return fmt.Errorf("compileIfStatement: symbol '{' is missing, got %v", c.tokenizer.Current)
		} else {
			statement.Tokens = append(statement.Tokens, *open)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return err
		}

		// statements
		statements, err := c.CompileStatements()
		if err != nil {
			return fmt.Errorf("compileIfStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, statements)

		// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()

		// '}'
		if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('}') {
			return fmt.Errorf("compileIfStatement: symbol '}' is missing, got %v", c.tokenizer.Current)
		} else {
			statement.Tokens = append(statement.Tokens, *close)
		}

		return nil
	}

	if err := consumeStatements(); err != nil {
		return nil, err
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// 'else'
	if kwd, err := c.tokenizer.Keyword(); err == nil && kwd.Val() == token.KwdElse {
		statement.Tokens = append(statement.Tokens, *kwd)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		if err := consumeStatements(); err != nil {
			return nil, err
		}

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	}

	return &statement, nil
}

func (c *Compiler) compileWhileStatement() (*WhileStatement, error) {
	statement := WhileStatement{Tokens: []token.Element{}}

	// 'while'
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdWhile {
		return nil, fmt.Errorf("compileWhileStatement: whileStatement should start with WHILE, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// '('
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('(') {
		return nil, fmt.Errorf("compileWhileStatement: symbol '(' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ')'
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune(')') {
		return nil, fmt.Errorf("compileWhileStatement: symbol ')' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// '{'
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('{') {
		return nil, fmt.Errorf("compileWhileStatement: symbol '{' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// statements
	statements, err := c.CompileStatements()
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, statements)

	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()

	// '}'
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune('}') {
		return nil, fmt.Errorf("compileWhileStatement: symbol '}' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *close)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	return &statement, nil
}

func (c *Compiler) compileDoStatement() (*DoStatement, error) {
	statement := DoStatement{Tokens: []token.Element{}}

	// 'do'
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdDo {
		return nil, fmt.Errorf("compileDoStatement: doStatement should start with DO, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// subroutineCall
	call, err := c.CompileSubroutineCall()
	if err != nil {
		return nil, fmt.Errorf("compileDoStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, call...)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ';'
	if end, err := c.tokenizer.Symbol(); err != nil || end.Val() != rune(';') {
		return nil, fmt.Errorf("compileDoStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *end)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	return &statement, nil
}

func (c *Compiler) compileReturnStatement() (*ReturnStatement, error) {
	statement := ReturnStatement{Tokens: []token.Element{}}

	// 'return'
	if kwd, err := c.tokenizer.Keyword(); err != nil || kwd.Val() != token.KwdReturn {
		return nil, fmt.Errorf("compileReturnStatement: ReturnStatement should start with RETURN, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *kwd)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// expression?
	if end, err := c.tokenizer.Symbol(); !(err == nil && end.Val() == rune(';')) {
		exp, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("compileReturnStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, exp)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	}

	// ';'
	if end, err := c.tokenizer.Symbol(); err != nil || end.Val() != rune(';') {
		return nil, fmt.Errorf("compileReturnStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, *end)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	return &statement, nil
}
