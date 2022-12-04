package token

import (
	"regexp"
	"strconv"
)

type TokenType int

const (
	TkKeyword TokenType = iota
	TkSymbol
	TkIdentifier
	TkIntConst
	TkStringConst
)

type Kwd int

func (kwd Kwd) String() string {
	label := map[Kwd]string{
		KwdClass:       "class",
		KwdConstructor: "constructor",
		KwdFunction:    "function",
		KwdMethod:      "method",
		KwdField:       "field",
		KwdStatic:      "static",
		KwdVar:         "var",
		KwdInt:         "int",
		KwdChar:        "char",
		KwdBoolean:     "boolean",
		KwdVoid:        "void",
		KwdTrue:        "true",
		KwdFalse:       "false",
		KwdNull:        "null",
		KwdThis:        "this",
		KwdLet:         "let",
		KwdDo:          "do",
		KwdIf:          "if",
		KwdElse:        "else",
		KwdWhile:       "while",
		KwdReturn:      "return",
	}
	return label[kwd]
}

const (
	KwdClass Kwd = iota
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

var keywords = []Kwd{
	KwdClass,
	KwdConstructor,
	KwdFunction,
	KwdMethod,
	KwdField,
	KwdStatic,
	KwdVar,
	KwdInt,
	KwdChar,
	KwdBoolean,
	KwdVoid,
	KwdTrue,
	KwdFalse,
	KwdNull,
	KwdThis,
	KwdLet,
	KwdDo,
	KwdIf,
	KwdElse,
	KwdWhile,
	KwdReturn,
}

func toKeyword(s string) *Kwd {
	for _, k := range keywords {
		if s == k.String() {
			return &k
		}
	}
	return nil
}

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

func toSymbol(r rune) *rune {
	for _, tk := range symbols {
		if r == tk {
			return &r
		}
	}
	return nil
}

const intConstMax = 32767

func toIntConst(s string) *int {
	if v, err := strconv.Atoi(s); err != nil {
		return nil
	} else {
		if v > intConstMax {
			return nil
		}
		return &v
	}
}

var strConstRegex = regexp.MustCompile(`^"(?P<val>[^"\n]*)"$`)

func toStrConst(s string) *string {
	matches := strConstRegex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return &matches[strConstRegex.SubexpIndex("val")]
	}
	return nil
}

var idRegex = regexp.MustCompile(`^(?P<val>[a-zA-Z_][0-9a-zA-Z_]*)$`)

func toID(s string) *string {
	matches := idRegex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return &matches[idRegex.SubexpIndex("val")]
	}
	return nil
}
