package token

type Token int

const (
	TokenKeyword Token = iota
	TokenSymbol
	TokenIdentifier
	TokenIntConst
	TokenStringConst
)

type Kwd int

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
