package token

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Tokenizer struct {
	TkType TokenType

	reader        *bufio.Reader
	hasMoreTokens bool

	keyword Keyword
	symbol  Symbol
	id      Identifier
	intVal  IntConst
	strVal  StringConst
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
	return t.TkType
}

func (t *Tokenizer) Keyword() (kwd *Keyword, err error) {
	if t.TkType != TkKeyword {
		err = fmt.Errorf("Keyword: token type is invalid: %d", t.TkType)
		return
	}

	kwd = &t.keyword
	return
}

func (t *Tokenizer) Symbol() (symbol *Symbol, err error) {
	if t.TkType != TkSymbol {
		err = fmt.Errorf("Symbol: token type is invalid: %d", t.TkType)
		return
	}
	symbol = &t.symbol
	return
}

func (t *Tokenizer) Identifier() (id *Identifier, err error) {
	if t.TkType != TkIdentifier {
		err = fmt.Errorf("Identifier: token type is invalid: %d", t.TkType)
		return
	}

	id = &t.id
	return
}

func (t *Tokenizer) IntVal() (v *IntConst, err error) {
	if t.TkType != TkIntConst {
		err = fmt.Errorf("IntVal: token type is invalid: %d", t.TkType)
		return
	}

	v = &t.intVal
	return
}

func (t *Tokenizer) StringVal() (v *StringConst, err error) {
	if t.TkType != TkStringConst {
		err = fmt.Errorf("StringVal: token type is invalid: %d", t.TkType)
		return
	}

	v = &t.strVal
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

func (t *Tokenizer) tokenize() (err error) {
	defer func() {
		if err == io.EOF {
			t.hasMoreTokens = false
			err = nil
		}
	}()

	r, _, err := t.reader.ReadRune()
	if err != nil {
		return
	}

	// 空白の場合、後に連続する空白をすべて読んでrerun
	if unicode.IsSpace(r) {
		t.TkType = TkWhiteSpace
		err = t.consumeWhiteSpaces()
		err = t.tokenize()
		return
	}

	// コメントの場合、コメントをすべて読んでrerun
	if r == rune('/') {
		next, _, e := t.reader.ReadRune()
		if e != nil {
			err = e
			return
		}

		// '//'または'/*'で始まる場合はコメントと判定
		if next == rune('/') {
			t.TkType = TkComment
			err = t.consumeInlineComment()
			err = t.tokenize()
			return
		}
		if next == rune('*') {
			t.TkType = TkComment
			err = t.consumeMultilineComment()
			err = t.tokenize()
			return
		}

		// コメントではない場合, 先読みした分をUnread
		if err = t.reader.UnreadRune(); err != nil {
			return
		}
	}

	// symbolかチェック
	symbol, err := toSymbol(r)
	if err != nil {
		return
	}
	if symbol != nil {
		t.TkType = TkSymbol
		t.symbol = *symbol
		return
	}

	// 一文字目が'"'の場合、文字列としてparse
	isStrConst := r == rune('"')

	runes := []rune{r}
	for {
		r, _, e := t.reader.ReadRune()
		if e != nil {
			err = e
			return
		}

		// 文字列の場合'"', それ以外の場合シンボルまたは空白が見つかったらbreak
		if isStrConst {
			if r == rune('"') {
				runes = append(runes, r)
				break
			}
		} else {
			symbol, e := toSymbol(r)
			if e != nil {
				err = e
				return
			}
			if symbol != nil || unicode.IsSpace(r) {
				if err = t.reader.UnreadRune(); err != nil {
					return
				}
				break
			}
		}

		runes = append(runes, r)
	}

	s := string(runes)

	if kwd := toKeyword(s); kwd != nil {
		t.TkType = TkKeyword
		t.keyword = *kwd
		return
	}

	if i := toIntConst(s); i != nil {
		t.TkType = TkIntConst
		t.intVal = *i
		return
	}

	if str := toStrConst(s); str != nil {
		t.TkType = TkStringConst
		t.strVal = *str
		return
	}

	if id := toID(s); id != nil {
		t.TkType = TkIdentifier
		t.id = *id
		return
	}

	err = fmt.Errorf("unexpected statement: %s", s)
	return
}
