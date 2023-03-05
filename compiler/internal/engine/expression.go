package engine

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/symtab"
	token "github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

type Expression struct {
	XMLName xml.Name `xml:"expression"`
	Tokens  []token.Element
}

func (el Expression) ElementType() token.ElementType {
	return token.ElExpression
}

type ExpressionList struct {
	XMLName xml.Name `xml:"expressionList"`
	Tokens  []token.Element
}

func (el ExpressionList) ElementType() token.ElementType {
	return token.ElExpressionList
}

type Term struct {
	XMLName xml.Name `xml:"term"`
	Tokens  []token.Element
}

func (el Term) ElementType() token.ElementType {
	return token.ElTerm
}

func (c *Compiler) CompileExpression() (*Expression, error) {
	expression := Expression{Tokens: []token.Element{}}

	// term (op term)*
	for {
		term, err := c.CompileTerm()
		if err != nil {
			return nil, fmt.Errorf("compileExpression: %w", err)
		}
		expression.Tokens = append(expression.Tokens, term)

		// check that the current token is op
		op, err := c.tokenizer.Symbol()
		if err != nil {
			break
		}
		v := op.Val()
		if !(v == token.SymPlus ||
			v == token.SymMinus ||
			v == token.SymAsterisk ||
			v == token.SymSlash ||
			v == token.SymAmpersand ||
			v == token.SymBar ||
			v == token.SymLessThan ||
			v == token.SymGreaterThan ||
			v == token.SymEqual) {
			break
		}
		expression.Tokens = append(expression.Tokens, op)
		if err := c.tokenizer.Advance(); err != nil {
			return nil, fmt.Errorf("compileExpression: %w", err)
		}
	}

	return &expression, nil
}

func (c *Compiler) CompileTerm() (*Term, error) {
	term := Term{Tokens: []token.Element{}}

	tkType := c.tokenizer.TokenType()

	switch tkType {
	// integerConstant
	case token.TkIntConst:
		// ignore the error because it is already checked that the token type is INT_CONST
		v, _ := c.tokenizer.IntVal()
		term.Tokens = append(term.Tokens, v)
		return &term, c.tokenizer.Advance()

	// stringConstant
	case token.TkStringConst:
		// ignore the error because it is already checked that the token type is STRING_CONST
		v, _ := c.tokenizer.StringVal()
		term.Tokens = append(term.Tokens, v)
		return &term, c.tokenizer.Advance()

	// keywordConstant
	case token.TkKeyword:
		kwd, err := c.consumeKeyword(token.KwdTrue, token.KwdFalse, token.KwdNull, token.KwdThis)
		if err != nil {
			return nil, fmt.Errorf("CompileTerm: invalid keyword %s", kwd)
		}
		term.Tokens = append(term.Tokens, kwd)

	// varName | varName '[' expression ']' | subroutineCall
	// subroutineCall: subroutineName '(' expressionList ')' | (className | varName) '.' subroutineName '(' expressionList ')'
	case token.TkIdentifier:
		next, err := c.tokenizer.CheckNextRune()
		if err != nil {
			return nil, fmt.Errorf("compileTerm: %w", err)
		}

		switch next {
		// varName '[' expression ']'
		case token.SymLeftSquareBracket:
			tokens, err := c.CompileArrayDec()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, tokens...)

		// subroutineCall
		case token.SymLeftParenthesis, token.SymDot:
			tokens, err := c.CompileSubroutineCall()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, tokens...)

		// varName
		default:
			name, err := c.compileName()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, name)
		}

	// '(' expression ')' | unaryOp term
	case token.TkSymbol:
		// ignore the error because it is already checked that the token type is SYMBOL
		s, _ := c.tokenizer.Symbol()

		switch s.Val() {
		// '(' expression ')'
		case token.SymLeftParenthesis:
			term.Tokens = append(term.Tokens, s)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}

			expression, err := c.CompileExpression()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, expression)

			close, err := c.consumeSymbol(token.SymRightParenthesis)
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, close)

		// unaryOp term
		case token.SymMinus, token.SymTilde:
			term.Tokens = append(term.Tokens, s)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}

			t, err := c.CompileTerm()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, t)

		default:
			return nil, fmt.Errorf("CompileTerm: invalid symbol %s", s)
		}

	default:
		return nil, fmt.Errorf("CompileTerm: invalid token %v", c.tokenizer.Current)
	}

	return &term, nil
}

func (c *Compiler) CompileArrayDec() ([]token.Element, error) {
	tokens := []token.Element{}

	// varName
	name, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileArrayDec: %w", err)
	}
	tokens = append(tokens, name)

	// '['
	if open, err := c.consumeSymbol(token.SymLeftSquareBracket); err != nil {
		return nil, fmt.Errorf("CompileArrayDec: symbol '[' is missing, got %v", c.tokenizer.Current)
	} else {
		tokens = append(tokens, open)
	}

	// expression
	expression, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("CompileArrayDec: %w", err)
	}
	tokens = append(tokens, expression)

	// ']'
	if close, err := c.consumeSymbol(token.SymRightSquareBracket); err != nil {
		return nil, fmt.Errorf("CompileArrayDec: symbol ']' is missing, got %v", c.tokenizer.Current)
	} else {
		tokens = append(tokens, close)
	}

	return tokens, nil
}

func (c *Compiler) CompileSubroutineCall() ([]token.Element, error) {
	tokens := []token.Element{}

	// subroutineName or (className | varName)
	name, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileSubroutineCall: %w", err)
	}
	tokens = append(tokens, name)

	consumeExpressionList := func() error {
		// '('
		if open, err := c.consumeSymbol(token.SymLeftParenthesis); err != nil {
			return fmt.Errorf("CompileSubroutineCall: symbol '(' is missing, got %v", c.tokenizer.Current)
		} else {
			tokens = append(tokens, open)
		}

		// expressionList
		// NOTE: You don't need to call Advance() because Advance() is already called inside CompileExpressionList()
		list, err := c.CompileExpressionList()
		if err != nil {
			return fmt.Errorf("compileSubroutineCall: %w", err)
		}
		tokens = append(tokens, list)

		// ')'
		if close, err := c.consumeSymbol(token.SymRightParenthesis); err != nil {
			return fmt.Errorf("CompileSubroutineCall: symbol ')' is missing, got %v", c.tokenizer.Current)
		} else {
			tokens = append(tokens, close)
		}

		return nil
	}

	// '(' or '.'
	s, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("compileSubroutineCall: %w", err)
	}
	switch s.Val() {
	case token.SymLeftParenthesis:
		name.IsDefined = false
		name.Category = symtab.SkSubroutine.String()

		// '(' expressionList ')'
		if err := consumeExpressionList(); err != nil {
			return nil, err
		}
	case token.SymDot:
		name.IsDefined = false
		name.Category = symtab.SkClass.String()

		// NOTE: このappendいる？
		tokens = append(tokens, s)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// subroutineName
		id, err := c.compileName()
		if err != nil {
			return nil, fmt.Errorf("compileSubroutineCall: %w", err)
		}
		tokens = append(tokens, id)
		id.IsDefined = false
		name.Category = symtab.SkClass.String()

		// '(' expressionList ')'
		if err := consumeExpressionList(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("compileSubroutineCall: '(' or '.' is expected, got %s", s)
	}

	return tokens, nil
}

func (c *Compiler) CompileExpressionList() (*ExpressionList, error) {
	expressionList := &ExpressionList{Tokens: []token.Element{}}

	// Return empty ParameterList when current token is ')'
	if s, err := c.tokenizer.Symbol(); err == nil && s.Val() == token.SymRightParenthesis {
		return expressionList, nil
	}

	// (expression (',' expression)*)
	for {
		// expression
		expression, err := c.CompileExpression()
		if err != nil {
			return nil, fmt.Errorf("CompileParameterList: %w", err)
		}
		expressionList.Tokens = append(expressionList.Tokens, expression)

		// check additional parameter
		if s, err := c.consumeSymbol(token.SymComma); err != nil {
			break
		} else {
			expressionList.Tokens = append(expressionList.Tokens, s)
		}
	}

	return expressionList, nil
}
