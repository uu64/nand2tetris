package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
)

type Cmd int

const (
	UNKNOWN Cmd = iota
	COMMENT
	EMPTY
	A_CMD
	C_CMD
	L_CMD
)

type Parser struct {
	scanner         *bufio.Scanner
	hasMoreCommands bool
	currentCmd      Cmd
	symbol          []byte
	dest            []byte
	comp            []byte
	jump            []byte
}

func New(f io.Reader) *Parser {
	s := bufio.NewScanner(f)
	return &Parser{scanner: s, hasMoreCommands: true}
}

func (p *Parser) HasMoreCommands() bool {
	return p.hasMoreCommands
}

func (p *Parser) Advance() {
	if !p.scanner.Scan() {
		err := p.scanner.Err()
		// io.EOF
		if err == nil {
			p.hasMoreCommands = false
		} else {
			log.Fatal("Parser.Advance: %w", err)
		}
	}
	p.parse(p.scanner.Bytes())
}

func (p *Parser) CommandType() Cmd {
	return p.currentCmd
}

func (p *Parser) Symbol() string {
	return string(p.symbol)
}

func (p *Parser) Dest() string {
	return string(p.dest)
}

func (p *Parser) Comp() string {
	return string(p.comp)
}

func (p *Parser) Jump() string {
	return string(p.jump)
}

var aCmdPtn = regexp.MustCompile(`^@(?P<symbol>\w+|[0-9]+)`)

var cCmdPtn = regexp.MustCompile(`^(?P<dest>null|[AMD]+)?=?(?P<comp>[AMD01&|+\-\!]+);?(?P<jump>null|JGT|JEQ|JGE|JLT|JNE|JLE|JMP)?`)

func (p *Parser) parse(row []byte) {
	b := bytes.TrimSpace(row)

	// skip empty row
	if len(b) == 0 {
		p.currentCmd = EMPTY
		return
	}

	// skip comment
	if bytes.HasPrefix(b, []byte("//")) {
		p.currentCmd = COMMENT
		return
	}

	switch b[0] {
	case '@':
		// A command
		matches := aCmdPtn.FindSubmatch(b)
		if len(matches) > 0 {
			p.symbol = matches[aCmdPtn.SubexpIndex("symbol")]
		}
		p.currentCmd = A_CMD
	case '(':
		// L command
		p.currentCmd = L_CMD
		fmt.Println("L command")
	default:
		// C command
		matches := cCmdPtn.FindSubmatch(b)
		if len(matches) > 0 {
			p.dest = matches[cCmdPtn.SubexpIndex("dest")]
			p.comp = matches[cCmdPtn.SubexpIndex("comp")]
			p.jump = matches[cCmdPtn.SubexpIndex("jump")]
		}
		p.currentCmd = C_CMD
	}
}
