package token

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
)

// token
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

// keyword
type KeywordType int

const (
	KwdClass KeywordType = iota
	KwdMethod
	KwdFunction
	KwdConstructor
	KwdInt
	KwdBoolean
	KwdChar
	KwdVoid
	KwdVar
	KwdStatic
	KwdField
	KwdLet
	KwdDo
	KwdIf
	KwdElse
	KwdWhile
	KwdReturn
	KwdTrue
	KwdFalse
	KwdNull
	KwdThis
)

var KwdLabelMap = map[string]KeywordType{
	"class":       KwdClass,
	"constructor": KwdConstructor,
	"function":    KwdFunction,
	"method":      KwdMethod,
	"field":       KwdField,
	"static":      KwdStatic,
	"var":         KwdVar,
	"int":         KwdInt,
	"char":        KwdChar,
	"boolean":     KwdBoolean,
	"void":        KwdVoid,
	"true":        KwdTrue,
	"false":       KwdFalse,
	"null":        KwdNull,
	"this":        KwdThis,
	"let":         KwdLet,
	"do":          KwdDo,
	"if":          KwdIf,
	"else":        KwdElse,
	"while":       KwdWhile,
	"return":      KwdReturn,
}

type Keyword struct {
	XMLName xml.Name `xml:"keyword"`
	Label   string   `xml:",chardata"`
}

func (tk Keyword) TokenType() TokenType {
	return TkKeyword
}

func (tk Keyword) Val() KeywordType {
	return KwdLabelMap[tk.Label]
}

func toKeyword(s string) *Keyword {
	if _, ok := KwdLabelMap[s]; ok {
		return &Keyword{Label: s}
	}
	return nil
}

// symbol
var symbols = []rune{
	rune('{'),
	rune('}'),
	rune('('),
	rune(')'),
	rune('['),
	rune(']'),
	rune('.'),
	rune(','),
	rune(';'),
	rune('+'),
	rune('-'),
	rune('*'),
	rune('/'),
	rune('&'),
	rune('|'),
	rune('<'),
	rune('>'),
	rune('='),
	rune('~'),
}

type Symbol struct {
	XMLName xml.Name `xml:"symbol"`
	Label   string   `xml:",chardata"`
}

func (tk Symbol) TokenType() TokenType {
	return TkSymbol
}

func (tk Symbol) Val() rune {
	// TODO: stringをxml unescapeしてruneにする処理を後で実装する
	return rune('-')
}

func toSymbol(r rune) (*Symbol, error) {
	for _, tk := range symbols {
		if r == tk {
			return &Symbol{Label: string(r)}, nil
		}
	}
	return nil, nil
}

// integer const
const intConstMax = 32767

type IntConst struct {
	XMLName xml.Name `xml:"integerConstant"`
	Label   string   `xml:",chardata"`
}

func (tk IntConst) TokenType() TokenType {
	return TkIntConst
}

func (tk IntConst) Val() (int, error) {
	v, err := strconv.Atoi(tk.Label)
	if err != nil {
		return -1, err
	}
	if v > intConstMax {
		return -1, fmt.Errorf("%s is over %d", tk.Label, intConstMax)
	}
	return v, nil
}

func toIntConst(s string) *IntConst {
	v, err := strconv.Atoi(s)
	if err != nil || v > intConstMax {
		return nil
	}
	return &IntConst{Label: s}
}

// string const
var strConstRegex = regexp.MustCompile(`^"(?P<val>[^"\n]*)"$`)

type StringConst struct {
	XMLName xml.Name `xml:"stringConstant"`
	Label   string   `xml:",chardata"`
}

func (tk StringConst) TokenType() TokenType {
	return TkStringConst
}

func (tk StringConst) Val() string {
	return tk.Label
}

func toStrConst(s string) *StringConst {
	matches := strConstRegex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return &StringConst{
			Label: matches[strConstRegex.SubexpIndex("val")],
		}
	}
	return nil
}

// identifier
var idRegex = regexp.MustCompile(`^(?P<val>[a-zA-Z_][0-9a-zA-Z_]*)$`)

type Identifier struct {
	XMLName xml.Name `xml:"identifier"`
	Label   string   `xml:",chardata"`
}

func (tk Identifier) TokenType() TokenType {
	return TkIdentifier
}

func (tk Identifier) Val() string {
	return tk.Label
}

func toID(s string) *Identifier {
	matches := idRegex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return &Identifier{
			Label: matches[idRegex.SubexpIndex("val")],
		}
	}
	return nil
}
