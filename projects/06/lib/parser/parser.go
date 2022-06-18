package parser

import (
	"bufio"
	"log"
	"os"
)

type Command int

const (
	A_COMMAND Command = iota
	C_COMMAND
	L_COMMAND
)

const delim = byte('\n')

type Parser struct {
	reader  *bufio.Reader
	current []byte
}

func New(f *os.File) *Parser {
	r := bufio.NewReader(f)
	return &Parser{reader: r}
}

func (p *Parser) HasMoreCommands() bool { return true }

func (p *Parser) Advance() {
	b, err := p.reader.ReadBytes(delim)
	if err != nil {
		// hasMoreCommands() = trueの場合のみ呼ばれるはずなので
		// io.EOFは発生しないはず
		log.Fatal("Parser.Advance: %w", err)
	}
	p.current = b
}

func (p *Parser) CommandType() Command { return A_COMMAND }

func (p *Parser) Symbol() string { return "" }

func (p *Parser) Dest() string { return "" }

func (p *Parser) Comp() string { return "" }

func (p *Parser) Jump() string { return "" }
