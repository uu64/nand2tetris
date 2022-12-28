package token

import "encoding/xml"

type TokenType int

const (
	TkKeyword TokenType = iota
	TkSymbol
	TkIdentifier
	TkIntConst
	TkStringConst
	TkEOF
	TkErr
)

type Token interface {
	TokenType() TokenType
}

type EOF struct{}

func (eof EOF) TokenType() TokenType {
	return TkEOF
}

type Err struct{}

func (err Err) TokenType() TokenType {
	return TkErr
}

type Tokens struct {
	XMLName xml.Name `xml:"tokens"`
	Tokens  []Token
}
