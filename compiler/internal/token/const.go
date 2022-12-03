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
	kwd{KwdClass, "class"},
	kwd{KwdConstructor, "constructor"},
	kwd{KwdFunction, "function"},
	kwd{KwdMethod, "method"},
	kwd{KwdField, "field"},
	kwd{KwdStatic, "static"},
	kwd{KwdVar, "var"},
	kwd{KwdInt, "int"},
	kwd{KwdChar, "char"},
	kwd{KwdBoolean, "boolean"},
	kwd{KwdVoid, "void"},
	kwd{KwdTrue, "true"},
	kwd{KwdFalse, "false"},
	kwd{KwdNull, "null"},
	kwd{KwdThis, "this"},
	kwd{KwdLet, "let"},
	kwd{KwdDo, "do"},
	kwd{KwdIf, "if"},
	kwd{KwdElse, "else"},
	kwd{KwdWhile, "while"},
	kwd{KwdReturn, "return"},
}

func isKeyword(s string) (bool, *kwd) {
	for _, k := range keywords {
		if s == k.label {
			return true, &k
		}
	}
	return false, nil
}
