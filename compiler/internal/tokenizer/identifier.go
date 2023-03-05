package tokenizer

import (
	"encoding/xml"
	"regexp"
)

// identifier
var idRegex = regexp.MustCompile(`^(?P<val>[a-zA-Z_][0-9a-zA-Z_]*)$`)

type Identifier struct {
	XMLName   xml.Name `xml:"identifier"`
	Label     string   `xml:",chardata"`
	Category  string   `xml:"-"`
	Kind      string   `xml:"-"`
	Index     int      `xml:"-"`
	IsDefined bool     `xml:"-"`
	// Label     string `xml:"value"`
	// Category  string `xml:"category"`
	// Kind      string `xml:"kind"`
	// Index     int    `xml:"idx"`
	// IsDefined bool   `xml:"defined"`
}

func (tk *Identifier) TokenType() TokenType {
	return TkIdentifier
}

func (tk *Identifier) ElementType() ElementType {
	return ElToken
}

func (tk *Identifier) Val() string {
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
