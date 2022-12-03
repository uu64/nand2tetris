package token

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Tokenizer struct {
	reader        *bufio.Reader
	hasMoreTokens bool
	tkType        TokenType

	keyword kwd
	symbol  rune
	id      string
	intVal  int
	strVal  string
}

func New(f io.Reader) *Tokenizer {
	r := bufio.NewReader(f)
	return &Tokenizer{
		reader:        r,
		hasMoreTokens: true,
	}
}

func (t *Tokenizer) HasMoreTokens() bool {
	return t.hasMoreTokens
}

func (t *Tokenizer) Advance() error {
	if !t.HasMoreTokens() {
		return fmt.Errorf("Advance: no tokens")
	}

	if err := t.tokenize(); err != nil {
		return fmt.Errorf("Advance: tokenize failed: %w", err)
	}

	return nil
}

func (t *Tokenizer) TokenType() TokenType {
	return t.tkType
}

func (t *Tokenizer) Keyword() (id KwdID, err error) {
	if t.tkType != TkKeyword {
		err = fmt.Errorf("Keyword: token type is invalid: %d", t.tkType)
		return
	}

	id = t.keyword.id
	return
}

func (t *Tokenizer) Symbol() (symbol string, err error) {
	if t.tkType != TkSymbol {
		err = fmt.Errorf("Symbol: token type is invalid: %d", t.tkType)
		return
	}

	symbol = string(t.symbol)
	return
}

func (t *Tokenizer) Identifier() (id string, err error) {
	if t.tkType != TkIdentifier {
		err = fmt.Errorf("Identifier: token type is invalid: %d", t.tkType)
		return
	}

	id = t.id
	return
}

func (t *Tokenizer) IntVal() (v int, err error) {
	if t.tkType != TkIntConst {
		err = fmt.Errorf("IntVal: token type is invalid: %d", t.tkType)
		return
	}

	v = t.intVal
	return
}

func (t *Tokenizer) StringVal() (v string, err error) {
	if t.tkType != TkStringConst {
		err = fmt.Errorf("StringVal: token type is invalid: %d", t.tkType)
		return
	}

	v = t.strVal
	return
}

func (t *Tokenizer) consumeWhiteSpace() (rune, error) {
	for {
		r, _, err := t.reader.ReadRune()
		if err != nil {
			return r, err
		}

		if !unicode.IsSpace(r) {
			return r, nil
		}
	}
}

func (t *Tokenizer) tokenize() error {
	next, err := t.consumeWhiteSpace()
	if err != nil {
		if err == io.EOF {
			t.hasMoreTokens = false
			return nil
		}
		return err
	}

	if f, symbol := isSymbol(next); f {
		fmt.Println("↓ is token")
		fmt.Println(string(next))

		t.tkType = TkSymbol
		t.symbol = *symbol
		return nil
	}

	runes := []rune{next}
	for {
		r, _, err := t.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				t.hasMoreTokens = false
				return nil
			}
			return err
		}

		if f, _ := isSymbol(r); f || unicode.IsSpace(r) {
			if err := t.reader.UnreadRune(); err != nil {
				return err
			}
			break
		}

		runes = append(runes, r)
	}

	s := string(runes)
	if f, kwd := isKeyword(s); f {
		fmt.Println("↓ is keyword")
		fmt.Println(s)

		t.tkType = TkKeyword
		t.keyword = *kwd
		return nil

	}

	fmt.Println("↓ is unknown")
	fmt.Println(s)
	return nil
}
