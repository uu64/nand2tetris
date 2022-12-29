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

func (tkType TokenType) String() string {
	switch tkType {
	case TkKeyword:
		return "keyword"
	case TkSymbol:
		return "symbol"
	case TkIdentifier:
		return "identifier"
	case TkIntConst:
		return "integerConstant"
	case TkStringConst:
		return "stringConstant"
	case TkEOF:
		return "EOF"
	case TkErr:
		return "Err"
	}
	panic("unknown token type")
}

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
