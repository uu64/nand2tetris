package engine

import (
	"encoding/xml"
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/symtab"
	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
)

type Expression struct {
	XMLName xml.Name `xml:"expression"`
	Tokens  []tokenizer.Element
}

func (el Expression) ElementType() tokenizer.ElementType {
	return tokenizer.ElExpression
}

type ExpressionList struct {
	XMLName xml.Name `xml:"expressionList"`
	Tokens  []tokenizer.Element
	Len     int `xml:"-"`
}

func (el ExpressionList) ElementType() tokenizer.ElementType {
	return tokenizer.ElExpressionList
}

type Term struct {
	XMLName xml.Name `xml:"term"`
	Tokens  []tokenizer.Element
}

func (el Term) ElementType() tokenizer.ElementType {
	return tokenizer.ElTerm
}

type Op tokenizer.Symbol

func (el Op) ElementType() tokenizer.ElementType {
	return tokenizer.ElOp
}

type UnaryOp tokenizer.Symbol

func (el UnaryOp) ElementType() tokenizer.ElementType {
	return tokenizer.ElUnaryOp
}

func (c *Compiler) CompileExpression() (*Expression, error) {
	expression := Expression{Tokens: []tokenizer.Element{}}
	ops := []Op{}

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
		if !(v == tokenizer.SymPlus ||
			v == tokenizer.SymMinus ||
			v == tokenizer.SymAsterisk ||
			v == tokenizer.SymSlash ||
			v == tokenizer.SymAmpersand ||
			v == tokenizer.SymBar ||
			v == tokenizer.SymLessThan ||
			v == tokenizer.SymGreaterThan ||
			v == tokenizer.SymEqual) {
			break
		}
		expression.Tokens = append(expression.Tokens, Op(*op))
		ops = append(ops, Op(*op))
		if err := c.tokenizer.Advance(); err != nil {
			return nil, fmt.Errorf("compileExpression: %w", err)
		}
	}

	for i := len(ops) - 1; i >= 0; i-- {
		c.writeOp(ops[i])
	}

	return &expression, nil
}

func (c *Compiler) CompileTerm() (*Term, error) {
	term := Term{Tokens: []tokenizer.Element{}}

	tkType := c.tokenizer.TokenType()

	switch tkType {
	// integerConstant
	case tokenizer.TkIntConst:
		// ignore the error because it is already checked that the token type is INT_CONST
		v, _ := c.tokenizer.IntVal()
		term.Tokens = append(term.Tokens, v)

		n, err := v.Val()
		if err != nil {
			return nil, err
		}
		c.writePushIntConst(n)

		return &term, c.tokenizer.Advance()

	// stringConstant
	case tokenizer.TkStringConst:
		// ignore the error because it is already checked that the token type is STRING_CONST
		v, _ := c.tokenizer.StringVal()
		term.Tokens = append(term.Tokens, v)
		return &term, c.tokenizer.Advance()

	// keywordConstant
	case tokenizer.TkKeyword:
		kwd, err := c.consumeKeyword(tokenizer.KwdTrue, tokenizer.KwdFalse, tokenizer.KwdNull, tokenizer.KwdThis)
		if err != nil {
			return nil, fmt.Errorf("CompileTerm: invalid keyword %s", kwd)
		}
		term.Tokens = append(term.Tokens, kwd)

	// varName | varName '[' expression ']' | subroutineCall
	// subroutineCall: subroutineName '(' expressionList ')' | (className | varName) '.' subroutineName '(' expressionList ')'
	case tokenizer.TkIdentifier:
		next, err := c.tokenizer.CheckNextRune()
		if err != nil {
			return nil, fmt.Errorf("compileTerm: %w", err)
		}

		switch next {
		// varName '[' expression ']'
		case tokenizer.SymLeftSquareBracket:
			tokens, err := c.CompileArrayDec()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, tokens...)

		// subroutineCall
		case tokenizer.SymLeftParenthesis, tokenizer.SymDot:
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
	case tokenizer.TkSymbol:
		// ignore the error because it is already checked that the token type is SYMBOL
		s, _ := c.tokenizer.Symbol()

		switch s.Val() {
		// '(' expression ')'
		case tokenizer.SymLeftParenthesis:
			term.Tokens = append(term.Tokens, s)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}

			expression, err := c.CompileExpression()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, expression)

			close, err := c.consumeSymbol(tokenizer.SymRightParenthesis)
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, close)

		// unaryOp term
		case tokenizer.SymMinus, tokenizer.SymTilde:
			op := UnaryOp(*s)
			term.Tokens = append(term.Tokens, op)
			if err := c.tokenizer.Advance(); err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}

			t, err := c.CompileTerm()
			if err != nil {
				return nil, fmt.Errorf("compileTerm: %w", err)
			}
			term.Tokens = append(term.Tokens, t)

			c.writeUnaryOp(op)
		default:
			return nil, fmt.Errorf("CompileTerm: invalid symbol %s", s)
		}

	default:
		return nil, fmt.Errorf("CompileTerm: invalid token %v", c.tokenizer.Current)
	}

	return &term, nil
}

func (c *Compiler) CompileArrayDec() ([]tokenizer.Element, error) {
	tokens := []tokenizer.Element{}

	// varName
	name, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileArrayDec: %w", err)
	}
	tokens = append(tokens, name)

	// '['
	if open, err := c.consumeSymbol(tokenizer.SymLeftSquareBracket); err != nil {
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
	if close, err := c.consumeSymbol(tokenizer.SymRightSquareBracket); err != nil {
		return nil, fmt.Errorf("CompileArrayDec: symbol ']' is missing, got %v", c.tokenizer.Current)
	} else {
		tokens = append(tokens, close)
	}

	return tokens, nil
}

func (c *Compiler) CompileSubroutineCall() ([]tokenizer.Element, error) {
	tokens := []tokenizer.Element{}

	// subroutineName or (className | varName)
	name, err := c.compileName()
	if err != nil {
		return nil, fmt.Errorf("compileSubroutineCall: %w", err)
	}
	tokens = append(tokens, name)

	consumeExpressionList := func() (*ExpressionList, error) {
		// '('
		if open, err := c.consumeSymbol(tokenizer.SymLeftParenthesis); err != nil {
			return nil, fmt.Errorf("CompileSubroutineCall: symbol '(' is missing, got %v", c.tokenizer.Current)
		} else {
			tokens = append(tokens, open)
		}

		// expressionList
		// NOTE: You don't need to call Advance() because Advance() is already called inside CompileExpressionList()
		list, err := c.CompileExpressionList()
		if err != nil {
			return nil, fmt.Errorf("compileSubroutineCall: %w", err)
		}
		tokens = append(tokens, list)

		// ')'
		if close, err := c.consumeSymbol(tokenizer.SymRightParenthesis); err != nil {
			return nil, fmt.Errorf("CompileSubroutineCall: symbol ')' is missing, got %v", c.tokenizer.Current)
		} else {
			tokens = append(tokens, close)
		}

		return list, nil
	}

	// '(' or '.'
	s, err := c.tokenizer.Symbol()
	if err != nil {
		return nil, fmt.Errorf("compileSubroutineCall: %w", err)
	}
	switch s.Val() {
	case tokenizer.SymLeftParenthesis:
		name.IsDefined = false
		name.Category = symtab.SkSubroutine.String()

		// '(' expressionList ')'
		list, err := consumeExpressionList()
		if err != nil {
			return nil, err
		}

		c.writeCall(name.Label, list.Len)
	case tokenizer.SymDot:
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
		list, err := consumeExpressionList()
		if err != nil {
			return nil, err
		}

		c.writeCall(fmt.Sprintf("%s.%s", name.Label, id.Label), list.Len)
	default:
		return nil, fmt.Errorf("compileSubroutineCall: '(' or '.' is expected, got %s", s)
	}

	return tokens, nil
}

func (c *Compiler) CompileExpressionList() (*ExpressionList, error) {
	expressionList := &ExpressionList{Tokens: []tokenizer.Element{}, Len: 0}

	// Return empty ParameterList when current token is ')'
	if s, err := c.tokenizer.Symbol(); err == nil && s.Val() == tokenizer.SymRightParenthesis {
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
		expressionList.Len += 1

		// check additional parameter
		if s, err := c.consumeSymbol(tokenizer.SymComma); err != nil {
			break
		} else {
			expressionList.Tokens = append(expressionList.Tokens, s)
		}
	}

	return expressionList, nil
}
