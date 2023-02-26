package tokenizer

import "encoding/xml"

type SymbolType rune

const (
	SymLeftCurlyBracket   = rune('{')
	SymLeftParenthesis    = rune('(')
	SymLeftSquareBracket  = rune('[')
	SymRightCurlyBracket  = rune('}')
	SymRightParenthesis   = rune(')')
	SymRightSquareBracket = rune(']')
	SymDot                = rune('.')
	SymComma              = rune(',')
	SymSemiColon          = rune(';')
	SymPlus               = rune('+')
	SymMinus              = rune('-')
	SymAsterisk           = rune('*')
	SymSlash              = rune('/')
	SymAmpersand          = rune('&')
	SymBar                = rune('|')
	SymLessThan           = rune('<')
	SymGreaterThan        = rune('>')
	SymEqual              = rune('=')
	SymTilde              = rune('~')
)

var symbols = []rune{
	SymLeftCurlyBracket,
	SymLeftParenthesis,
	SymLeftSquareBracket,
	SymRightCurlyBracket,
	SymRightParenthesis,
	SymRightSquareBracket,
	SymDot,
	SymComma,
	SymSemiColon,
	SymPlus,
	SymMinus,
	SymAsterisk,
	SymSlash,
	SymAmpersand,
	SymBar,
	SymLessThan,
	SymGreaterThan,
	SymEqual,
	SymTilde,
}

type Symbol struct {
	XMLName xml.Name `xml:"symbol"`
	Label   string   `xml:",chardata"`
	val     rune     `xml:"-"`
}

func (tk *Symbol) TokenType() TokenType {
	return TkSymbol
}

func (tk *Symbol) ElementType() ElementType {
	return ElToken
}

func (tk *Symbol) Val() rune {
	return tk.val
}

func (tk *Symbol) String() string {
	return string(tk.val)
}

func toSymbol(r rune) (*Symbol, error) {
	for _, tk := range symbols {
		if r == tk {
			return &Symbol{
				Label: string(r),
				val:   r,
			}, nil
		}
	}
	return nil, nil
}
