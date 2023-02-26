package tokenizer

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

const intConstMax = 32767

type IntConst struct {
	XMLName xml.Name `xml:"integerConstant"`
	Label   string   `xml:",chardata"`
}

func (tk *IntConst) TokenType() TokenType {
	return TkIntConst
}

func (tk *IntConst) ElementType() ElementType {
	return ElToken
}

func (tk *IntConst) Val() (int, error) {
	v, err := strconv.Atoi(tk.Label)
	if err != nil {
		return -1, err
	}
	if v > intConstMax {
		return -1, fmt.Errorf("%s is over %d", tk.Label, intConstMax)
	}
	return v, nil
}

func toIntConst(s string) *IntConst {
	v, err := strconv.Atoi(s)
	if err != nil || v > intConstMax {
		return nil
	}
	return &IntConst{Label: s}
}
