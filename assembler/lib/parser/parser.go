package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Command int

const (
	A_COMMAND Command = iota
	C_COMMAND
	L_COMMAND
)

var aCommandPattern = regexp.MustCompile(`^@(?P<value>\w+|[0-9]+)$`)

type Parser struct {
	scanner         *bufio.Scanner
	hasMoreCommands bool
	currentCmd      Command
}

func New(f *os.File) *Parser {
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

func (p *Parser) CommandType() Command {
	return p.currentCmd
}

func (p *Parser) Symbol() string {
	return string(p.scanner.Bytes())
}

func (p *Parser) Dest() string { return "" }

func (p *Parser) Comp() string { return "" }

func (p *Parser) Jump() string { return "" }

func (p *Parser) parse(row []byte) {
	b := bytes.TrimSpace(row)

	// skip empty row
	if len(b) == 0 {
		return
	}

	// skip comment
	if bytes.HasPrefix(b, []byte("//")) {
		return
	}

	fmt.Println(string(b))
	switch b[0] {
	case '@':
		matches := aCommandPattern.FindSubmatch(b)
		fmt.Println(string(matches[aCommandPattern.SubexpIndex("value")]))
		if aCommandPattern.Match(b) {
			fmt.Println("match")
		}
		p.currentCmd = A_COMMAND
	case '(':
		// L command
		p.currentCmd = L_COMMAND
	default:
		fmt.Println("C command")
		p.currentCmd = C_COMMAND
	}
}
