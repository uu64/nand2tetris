package compile

import (
	"encoding/xml"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Expression struct {
	XMLName xml.Name `xml:"expression"`
	Tokens  []token.Element
}

func (el Expression) ElementType() token.ElementType {
	return token.ElExpression
}

func (c *Compiler) CompileExpression() (*Expression, error) {
	expression := Expression{Tokens: []token.Element{}}

	return &expression, nil
}

func (c *Compiler) CompileSubroutineCall() ([]token.Element, error) {
	tokens := []token.Element{}
	return tokens, nil

}
