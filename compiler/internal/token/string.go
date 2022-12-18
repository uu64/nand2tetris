package token

import (
	"encoding/xml"
	"regexp"
)

var strConstRegex = regexp.MustCompile(`^"(?P<val>[^"\n]*)"$`)

type StringConst struct {
	XMLName xml.Name `xml:"stringConstant"`
	Label   string   `xml:",chardata"`
}

func (tk StringConst) TokenType() TokenType {
	return TkStringConst
}

func (tk StringConst) Val() string {
	return tk.Label
}

func toStrConst(s string) *StringConst {
	matches := strConstRegex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return &StringConst{
			Label: matches[strConstRegex.SubexpIndex("val")],
		}
	}
	return nil
}
