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
