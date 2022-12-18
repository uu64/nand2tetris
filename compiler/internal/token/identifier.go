package token

import (
	"encoding/xml"
	"regexp"
)

// identifier
var idRegex = regexp.MustCompile(`^(?P<val>[a-zA-Z_][0-9a-zA-Z_]*)$`)

type Identifier struct {
	XMLName xml.Name `xml:"identifier"`
	Label   string   `xml:",chardata"`
}

func (tk Identifier) TokenType() TokenType {
	return TkIdentifier
}

func (tk Identifier) Val() string {
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
