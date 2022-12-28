package compile

import (
	"encoding/xml"

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
