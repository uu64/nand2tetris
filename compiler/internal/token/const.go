package token

type TokenType int

const (
	TkKeyword TokenType = iota
	TkSymbol
	TkIdentifier
	TkIntConst
	TkStringConst
)

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

func isSymbol(r rune) (bool, *rune) {
	for _, tk := range symbols {
		if r == tk {
			return true, &r
		}
	}
	return false, nil
}

type KwdID int

const (
	KwdClass KwdID = iota
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

type kwd struct {
	id    KwdID
	label string
}

var keywords = []kwd{
	{KwdClass, "class"},
	{KwdConstructor, "constructor"},
	{KwdFunction, "function"},
	{KwdMethod, "method"},
	{KwdField, "field"},
	{KwdStatic, "static"},
	{KwdVar, "var"},
	{KwdInt, "int"},
	{KwdChar, "char"},
	{KwdBoolean, "boolean"},
	{KwdVoid, "void"},
	{KwdTrue, "true"},
	{KwdFalse, "false"},
	{KwdNull, "null"},
	{KwdThis, "this"},
	{KwdLet, "let"},
	{KwdDo, "do"},
	{KwdIf, "if"},
	{KwdElse, "else"},
	{KwdWhile, "while"},
	{KwdReturn, "return"},
}

func isKeyword(s string) (bool, *kwd) {
	for _, k := range keywords {
		if s == k.label {
			return true, &k
		}
	}
	return false, nil
}
