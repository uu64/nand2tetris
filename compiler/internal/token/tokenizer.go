package token

import (
	"bufio"
	"io"
)

type Tokenizer struct {
	scanner *bufio.Scanner
}

func New(f io.Reader) *Tokenizer {
	s := bufio.NewScanner(f)
	return &Tokenizer{s}
}

func (t *Tokenizer) HasMoreTokens() bool { return true }
func (t *Tokenizer) Advance() bool       { return true }
func (t *Tokenizer) TokenType() Token    { return TokenKeyword }
func (t *Tokenizer) Keyword() Kwd        { return KwdClass }
func (t *Tokenizer) Symbol() bool        { return true }
func (t *Tokenizer) Identifier() string  { return "" }
func (t *Tokenizer) IntVal() int         { return 0 }
func (t *Tokenizer) StringVal() string   { return "" }
