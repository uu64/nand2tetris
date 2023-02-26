package tokenizer

import (
	"encoding/xml"
	"regexp"
)

// identifier
var idRegex = regexp.MustCompile(`^(?P<val>[a-zA-Z_][0-9a-zA-Z_]*)$`)

type Identifier struct {
	XMLName xml.Name `xml:"identifier"`
	Label   string   `xml:",chardata"`
	// Label    string   `xml:"value"`
	// Category string   `xml:"category"`
	// Defined  bool     `xml:"defined"`
	// Kind     string   `xml:"kind"`
	// Index    int      `xml:"idx"`
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

// func (tk *Identifier) SetMetadata(category, kind string, index int, defined bool) {
// 	tk.Category = category
// 	tk.Kind = kind
// 	tk.Index = index
// 	tk.Defined = defined
// }

func toID(s string) *Identifier {
	matches := idRegex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return &Identifier{
			Label: matches[idRegex.SubexpIndex("val")],
		}
	}
	return nil
}
