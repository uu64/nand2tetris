package compile

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
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

	term, err := c.CompileTerm()
	if err != nil {
		return nil, fmt.Errorf("compileExpression: %w", err)
	}
	expression.Tokens = append(expression.Tokens, term)

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
	// stringConstant
	case token.TkStringConst:
		// ignore the error because it is already checked that the token type is STRING_CONST
		v, _ := c.tokenizer.StringVal()
		term.Tokens = append(term.Tokens, v)
	// keywordConstant
	case token.TkKeyword:
		// ignore the error because it is already checked that the token type is KEYWORD
		kwd, _ := c.tokenizer.Keyword()
		if v := kwd.Val(); v == token.KwdTrue || v == token.KwdFalse || v == token.KwdNull || v == token.KwdThis {
			term.Tokens = append(term.Tokens, kwd)
		} else {
			return nil, fmt.Errorf("CompileTerm: invalid keyword %s", kwd)
		}
	// varName | varName '[' expression ']' | subroutineCall
	// subroutineCall: subroutineName '(' expressionList ')' | (className | varName) '.' subroutineName '(' expressionList ')'
	case token.TkIdentifier:
		next, err := c.tokenizer.CheckNextRune()
		if err != nil {
			return nil, fmt.Errorf("compileTerm: %w", err)
		}

		switch next {
		// varName '[' expression ']'
		case rune('['):
			tokens, err := c.CompileArrayDec()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, tokens...)
		// subroutineCall
		case rune('('), rune('.'):
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
		case rune('('):
			return nil, fmt.Errorf("CompileTerm: '(' expression ')' not implemented")
		// unaryOp term
		case rune('-'), rune('~'):
			return nil, fmt.Errorf("CompileTerm: unaryOp term not implemented")
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
	name, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("compileArrayDec: %w", err)
	}
	tokens = append(tokens, name)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// '['
	if open, err := c.tokenizer.Symbol(); err != nil || open.Val() != rune('[') {
		return nil, fmt.Errorf("CompileArrayDec: symbol '[' is missing, got %v", c.tokenizer.Current)
	} else {
		tokens = append(tokens, *open)
	}

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// expression
	expression, err := c.CompileExpression()
	if err != nil {
		return nil, fmt.Errorf("CompileArrayDec: %w", err)
	}
	tokens = append(tokens, expression)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	// ']'
	if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune(']') {
		return nil, fmt.Errorf("CompileArrayDec: symbol ']' is missing, got %v", c.tokenizer.Current)
	} else {
		tokens = append(tokens, *close)
	}

	return tokens, nil
}

func (c *Compiler) CompileSubroutineCall() ([]token.Element, error) {
	tokens := []token.Element{}

	// subroutineName or (className | varName)
	name, err := c.tokenizer.Identifier()
	if err != nil {
		return nil, fmt.Errorf("compileSubroutineCall: %w", err)
	}
	tokens = append(tokens, name)

	if err := c.tokenizer.Advance(); err != nil {
		return nil, err
	}

	consumeExpressionList := func() error {
		// '('
		if open, err := c.tokenizer.Symbol(); err != nil || open.Val() != rune('(') {
			return fmt.Errorf("CompileSubroutineCall: symbol '(' is missing, got %v", c.tokenizer.Current)
		} else {
			tokens = append(tokens, *open)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return err
		}

		// expressionList
		list, err := c.CompileExpressionList()
		if err != nil {
			return fmt.Errorf("compileSubroutineCall: %w", err)
		}
		tokens = append(tokens, list)

		// NOTE: You don't need to call Advance() because Advance() is already called inside CompileExpressionList()

		// ')'
		if close, err := c.tokenizer.Symbol(); err != nil || close.Val() != rune(')') {
			return fmt.Errorf("CompileSubroutineCall: symbol ')' is missing, got %v", c.tokenizer.Current)
		} else {
			tokens = append(tokens, *close)
		}

		return nil
	}

	// '(' or '.'
	s, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("compileSubroutineCall: %w", err)
	}
	switch s.Val() {
	case rune('('):
		// '(' expressionList ')'
		if err := consumeExpressionList(); err != nil {
			return nil, err
		}
	case rune('.'):
		tokens = append(tokens, s)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// subroutineName
		id, err := c.tokenizer.Identifier()
		if err != nil {
			return nil, fmt.Errorf("compileSubroutineCall: %w", err)
		}
		tokens = append(tokens, id)

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

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
	if s, err := c.tokenizer.Symbol(); err == nil && s.Val() == rune(')') {
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

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}

		// check additional parameter
		if s, err := c.tokenizer.Symbol(); err != nil || s.Val() != rune(',') {
			break
		} else {
			expressionList.Tokens = append(expressionList.Tokens, s)
		}

		if err := c.tokenizer.Advance(); err != nil {
			return nil, err
		}
	}

	return expressionList, nil
}

func (c *Compiler) compileOp() (*token.Symbol, error) {
	s, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("compileOp: %w", err)
	}

	switch s.Val() {
	case rune('+'), rune('-'), rune('*'), rune('/'), rune('&'), rune('|'), rune('<'), rune('>'), rune('='):
		return s, nil
	default:
		return nil, fmt.Errorf("compileOp: invalid symbol %s", s)
	}
}
