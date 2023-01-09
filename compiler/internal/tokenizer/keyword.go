package tokenizer

import "encoding/xml"

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

func (tk Keyword) ElementType() ElementType {
	return ElToken
}

func (tk Keyword) Val() KeywordType {
	return KwdLabelMap[tk.Label]
}

func (tk Keyword) String() string {
	return tk.Label
}

func toKeyword(s string) *Keyword {
	if _, ok := KwdLabelMap[s]; ok {
		return &Keyword{Label: s}
	}
	return nil
}
