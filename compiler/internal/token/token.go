package token

import "encoding/xml"

type TokenType int

const (
	TkKeyword TokenType = iota
	TkSymbol
	TkIdentifier
	TkIntConst
	TkStringConst
	TkWhiteSpace
	TkComment
)

type Token interface {
	TokenType() TokenType
}

type Tokens struct {
	XMLName xml.Name `xml:"tokens"`
	Tokens  []Token
}
