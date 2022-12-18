package token

import "encoding/xml"

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
	val     rune     `xml:"-"`
}

func (tk Symbol) TokenType() TokenType {
	return TkSymbol
}

func (tk Symbol) Val() rune {
	return tk.val
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
