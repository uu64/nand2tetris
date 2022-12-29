package token

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Tokenizer struct {
	Current Token

	reader        *bufio.Reader
	hasMoreTokens bool
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
	return t.Current.TokenType()
}

func (t *Tokenizer) Keyword() (ptr *Keyword, err error) {
	if tkType := t.Current.TokenType(); tkType != TkKeyword {
		err = fmt.Errorf("Keyword: token type is invalid: %v", tkType)
		return
	}
	kwd := t.Current.(Keyword)
	ptr = &kwd
	return
}

func (t *Tokenizer) Symbol() (ptr *Symbol, err error) {
	if tkType := t.Current.TokenType(); tkType != TkSymbol {
		err = fmt.Errorf("Symbol: token type is invalid: %v", tkType)
		return
	}
	symbol := t.Current.(Symbol)
	ptr = &symbol
	return
}

func (t *Tokenizer) Identifier() (ptr *Identifier, err error) {
	if tkType := t.Current.TokenType(); tkType != TkIdentifier {
		err = fmt.Errorf("Identifier: token type is invalid: %v", tkType)
		return
	}
	id := t.Current.(Identifier)
	ptr = &id
	return
}

func (t *Tokenizer) IntVal() (ptr *IntConst, err error) {
	if tkType := t.Current.TokenType(); tkType != TkIntConst {
		err = fmt.Errorf("IntVal: token type is invalid: %v", tkType)
		return
	}
	v := t.Current.(IntConst)
	ptr = &v
	return
}

func (t *Tokenizer) StringVal() (ptr *StringConst, err error) {
	if tkType := t.Current.TokenType(); tkType != TkStringConst {
		err = fmt.Errorf("StringVal: token type is invalid: %v", tkType)
		return
	}
	v := t.Current.(StringConst)
	ptr = &v
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
			t.Current = EOF{}
			t.hasMoreTokens = false
			err = nil
		} else if err != nil {
			t.Current = Err{}
		}
	}()

	r, _, e := t.reader.ReadRune()
	if e != nil {
		err = e
		return
	}

	// 空白の場合、後に連続する空白をすべて読んでrerun
	if unicode.IsSpace(r) {
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
			err = t.consumeInlineComment()
			err = t.tokenize()
			return
		}
		if next == rune('*') {
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
		t.Current = *symbol
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
		t.Current = *kwd
		return
	}

	if i := toIntConst(s); i != nil {
		t.Current = *i
		return
	}

	if str := toStrConst(s); str != nil {
		t.Current = *str
		return
	}

	if id := toID(s); id != nil {
		t.Current = *id
		return
	}

	err = fmt.Errorf("unexpected statement: %s", s)
	return
}
