package compile

import (
	"encoding/xml"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Class struct {
	XMLName xml.Name `xml:"class"`
	Tokens  []token.Token
}
