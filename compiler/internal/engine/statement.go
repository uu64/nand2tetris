package engine

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
	"github.com/uu64/nand2tetris/compiler/internal/vmwriter"
)

type Statements struct {
	XMLName xml.Name `xml:"statements"`
	Tokens  []tokenizer.Element
}

func (el Statements) ElementType() tokenizer.ElementType {
	return tokenizer.ElStatements
}

type LetStatement struct {
	XMLName xml.Name `xml:"letStatement"`
	Tokens  []tokenizer.Element
}

func (el LetStatement) ElementType() tokenizer.ElementType {
	return tokenizer.ElLetStatement
}

type IfStatement struct {
	XMLName xml.Name `xml:"ifStatement"`
	Tokens  []tokenizer.Element
}

func (el IfStatement) ElementType() tokenizer.ElementType {
	return tokenizer.ElIfStatement
}

type WhileStatement struct {
	XMLName xml.Name `xml:"whileStatement"`
	Tokens  []tokenizer.Element
}

func (el WhileStatement) ElementType() tokenizer.ElementType {
	return tokenizer.ElWhileStatement
}

type DoStatement struct {
	XMLName xml.Name `xml:"doStatement"`
	Tokens  []tokenizer.Element
}

func (el DoStatement) ElementType() tokenizer.ElementType {
	return tokenizer.ElDoStatement
}

type ReturnStatement struct {
	XMLName xml.Name `xml:"returnStatement"`
	Tokens  []tokenizer.Element
}

func (el ReturnStatement) ElementType() tokenizer.ElementType {
	return tokenizer.ElReturnStatement
}

func (c *Compiler) CompileStatements() (*Statements, error) {
	statements := Statements{Tokens: []tokenizer.Element{}}

	for c.tokenizer.TokenType() == tokenizer.TkKeyword {
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()

		switch kwd.Val() {
		case tokenizer.KwdLet:
			statement, err := c.compileLetStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case tokenizer.KwdIf:
			statement, err := c.compileIfStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case tokenizer.KwdWhile:
			statement, err := c.compileWhileStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case tokenizer.KwdDo:
			statement, err := c.compileDoStatement()
			if err != nil {
				return nil, fmt.Errorf("CompileStatements: %w", err)
			}
			statements.Tokens = append(statements.Tokens, statement)
		case tokenizer.KwdReturn:
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
	statement := LetStatement{Tokens: []tokenizer.Element{}}

	// 'let'
	kwd, err := c.consumeKeyword(tokenizer.KwdLet)
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: letStatement should start with LET, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, kwd)

	// varName
	varName, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, varName)

	isArray := false
	// ('[' expression ']')?
	if open, err := c.consumeSymbol(tokenizer.SymLeftSquareBracket); err == nil {
		isArray = true

		// '['
		statement.Tokens = append(statement.Tokens, open)

		// expression
		exp, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("compileLetStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, exp)

		// ']'
		close, err := c.consumeSymbol(tokenizer.SymRightSquareBracket)
		if err != nil {
			return nil, fmt.Errorf("compileLetStatement: symbol ']' is missing, got %v", c.tokenizer.Current)
		}
		statement.Tokens = append(statement.Tokens, close)

		// write vm code for array
		c.writePushVar(*varName)
		c.codewriter.WriteArithmetic(vmwriter.Add)
	}

	// '='
	eq, err := c.consumeSymbol(tokenizer.SymEqual)
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: symbol '=' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, eq)

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	// ';'
	end, err := c.consumeSymbol(tokenizer.SymSemiColon)
	if err != nil {
		return nil, fmt.Errorf("compileLetStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, end)

	if isArray {
		// save the result of the expression
		c.writePopTemp(0)
		// save the memory address of the reference to the array
		c.writePopPointer(1)
		// restore the saved result and push it onto the stack
		c.writePushTemp(0)
		// assign a value to the memory addredss of the reference to array
		c.codewriter.WritePop(vmwriter.That, 0)
	} else {
		if err := c.writePopVar(*varName); err != nil {
			return nil, fmt.Errorf("compileLetStatement: %w", err)
		}
	}

	return &statement, nil
}

func (c *Compiler) compileIfStatement() (*IfStatement, error) {
	statement := IfStatement{Tokens: []tokenizer.Element{}}

	trueLabel := fmt.Sprintf("IF_TRUE%d", c.ctx.IfIndex)
	falseLabel := fmt.Sprintf("IF_FALSE%d", c.ctx.IfIndex)
	endLabel := fmt.Sprintf("IF_END%d", c.ctx.IfIndex)
	c.ctx.IfIndex += 1

	// 'if'
	kwd, err := c.consumeKeyword(tokenizer.KwdIf)
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: ifStatement should start with IF, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, kwd)

	// '('
	open, err := c.consumeSymbol(tokenizer.SymLeftParenthesis)
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: symbol '(' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, open)

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	// ')'
	close, err := c.consumeSymbol(tokenizer.SymRightParenthesis)
	if err != nil {
		return nil, fmt.Errorf("compileIfStatement: symbol ')' is missing, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, close)

	c.codewriter.WriteIf(trueLabel)
	c.codewriter.WriteGoTo(falseLabel)
	c.codewriter.WriteLabel(trueLabel)

	consumeStatements := func() error {
		// '{'
		open, err := c.consumeSymbol(tokenizer.SymLeftCurlyBracket)
		if err != nil {
			return fmt.Errorf("compileIfStatement: symbol '{' is missing, got %v", c.tokenizer.Current)
		}
		statement.Tokens = append(statement.Tokens, open)

		// statements
		// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()
		statements, err := c.CompileStatements()
		if err != nil {
			return fmt.Errorf("compileIfStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, statements)

		// '}'
		close, err := c.consumeSymbol(tokenizer.SymRightCurlyBracket)
		if err != nil {
			return fmt.Errorf("compileIfStatement: symbol '}' is missing, got %v", c.tokenizer.Current)
		}
		statement.Tokens = append(statement.Tokens, close)
		return nil
	}

	if err := consumeStatements(); err != nil {
		return nil, err
	}

	// 'else'
	if kwd, err := c.consumeKeyword(tokenizer.KwdElse); err == nil {
		c.codewriter.WriteGoTo(endLabel)
		c.codewriter.WriteLabel(falseLabel)

		statement.Tokens = append(statement.Tokens, kwd)
		if err := consumeStatements(); err != nil {
			return nil, err
		}

		c.codewriter.WriteLabel(endLabel)
	} else {
		c.codewriter.WriteLabel(falseLabel)
	}

	return &statement, nil
}

func (c *Compiler) compileWhileStatement() (*WhileStatement, error) {
	statement := WhileStatement{Tokens: []tokenizer.Element{}}

	startLabel := fmt.Sprintf("WHILE_EXP%d", c.ctx.WhileIndex)
	endLabel := fmt.Sprintf("WHILE_END%d", c.ctx.WhileIndex)
	c.ctx.WhileIndex += 1

	// 'while'
	kwd, err := c.consumeKeyword(tokenizer.KwdWhile)
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: whileStatement should start with WHILE, got %v", c.tokenizer.Current)
	}
	statement.Tokens = append(statement.Tokens, kwd)
	c.codewriter.WriteLabel(startLabel)

	// '('
	if open, err := c.consumeSymbol(tokenizer.SymLeftParenthesis); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol '(' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, open)
	}

	// expression
	exp, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, exp)

	// ')'
	if close, err := c.consumeSymbol(tokenizer.SymRightParenthesis); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol ')' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, close)
	}

	c.codewriter.WriteArithmetic(vmwriter.Not)
	c.codewriter.WriteIf(endLabel)

	// '{'
	if open, err := c.consumeSymbol(tokenizer.SymLeftCurlyBracket); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol '{' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, open)
	}

	// statements
	// NOTE: You don't need to call Advance() because Advance() is already called inside CompileStatements()
	statements, err := c.CompileStatements()
	if err != nil {
		return nil, fmt.Errorf("compileWhileStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, statements)

	// '}'
	if close, err := c.consumeSymbol(tokenizer.SymRightCurlyBracket); err != nil {
		return nil, fmt.Errorf("compileWhileStatement: symbol '}' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, close)
	}

	c.codewriter.WriteGoTo(startLabel)
	c.codewriter.WriteLabel(endLabel)

	return &statement, nil
}

func (c *Compiler) compileDoStatement() (*DoStatement, error) {
	statement := DoStatement{Tokens: []tokenizer.Element{}}

	// 'do'
	if kwd, err := c.consumeKeyword(tokenizer.KwdDo); err != nil {
		return nil, fmt.Errorf("compileDoStatement: doStatement should start with DO, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, kwd)
	}

	// subroutineCall
	call, err := c.CompileSubroutineCall()
	if err != nil {
		return nil, fmt.Errorf("compileDoStatement: %w", err)
	}
	statement.Tokens = append(statement.Tokens, call...)
	c.discardReturn()

	// ';'
	if end, err := c.consumeSymbol(tokenizer.SymSemiColon); err != nil {
		return nil, fmt.Errorf("compileDoStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, end)
	}

	return &statement, nil
}

func (c *Compiler) compileReturnStatement() (*ReturnStatement, error) {
	statement := ReturnStatement{Tokens: []tokenizer.Element{}}

	// 'return'
	if kwd, err := c.consumeKeyword(tokenizer.KwdReturn); err != nil {
		return nil, fmt.Errorf("compileReturnStatement: ReturnStatement should start with RETURN, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, kwd)
	}

	// expression?
	if end, err := c.tokenizer.Symbol(); !(err == nil && end.Val() == tokenizer.SymSemiColon) {
		exp, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("compileReturnStatement: %w", err)
		}
		statement.Tokens = append(statement.Tokens, exp)
	}

	// ';'
	if end, err := c.consumeSymbol(tokenizer.SymSemiColon); err != nil {
		return nil, fmt.Errorf("compileReturnStatement: symbol ';' is missing, got %v", c.tokenizer.Current)
	} else {
		statement.Tokens = append(statement.Tokens, end)
	}

	c.writeReturn()

	return &statement, nil
}
