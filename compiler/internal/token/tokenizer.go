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

func (t *Tokenizer) consumeWhiteSpaces() error {
	for {
		next, _, err := t.reader.ReadRune()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(next) {
			break
		}
	}
	// 連続した空白の次の最初の一文字を読んでいるので戻す
	return t.reader.UnreadRune()
}

func (t *Tokenizer) consumeInlineComment() error {
	// 行末まで読む
	if _, err := t.reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

func (t *Tokenizer) consumeMultilineComment() error {
	// '*/'まで読む
	for {
		if _, err := t.reader.ReadString('*'); err != nil {
			return err
		}

		r, _, err := t.reader.ReadRune()
		if err != nil {
			return err
		}

		if r == rune('/') {
			break
		}
	}
	return nil
}

func (t *Tokenizer) tokenize() error {
	r, _, err := t.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			t.hasMoreTokens = false
			return nil
		}
		return err
	}

	// 空白の場合、後に連続する空白をすべて読んでreturn
	if unicode.IsSpace(r) {
		return t.consumeWhiteSpaces()
	}

	// コメントの場合、コメントをすべて読んでreturn
	if r == rune('/') {
		next, _, err := t.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				t.hasMoreTokens = false
				return nil
			}
			return err
		}

		// '//'または'/*'で始まる場合はコメントと判定
		if next == rune('/') {
			return t.consumeInlineComment()
		}
		if next == rune('*') {
			return t.consumeMultilineComment()
		}

		// 先読みした分をUnread
		if err := t.reader.UnreadRune(); err != nil {
			return err
		}
	}

	// symbolかチェック
	if f, symbol := isSymbol(r); f {
		fmt.Printf("symbol: %s\n", string(r))

		t.tkType = TkSymbol
		t.symbol = *symbol
		return nil
	}

	runes := []rune{r}
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
		fmt.Printf("keyword: %s\n", s)

		t.tkType = TkKeyword
		t.keyword = *kwd
		return nil
	}

	fmt.Printf("unknown: %s\n", s)
	return nil
}
