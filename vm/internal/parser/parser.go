package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

type Cmd int

const (
	COMMENT Cmd = iota
	EMPTY
	C_ARITHMETRIC
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_CALL
	C_RETURN
)

type Parser struct {
	scanner         *bufio.Scanner
	hasMoreCommands bool
	currentCmd      Cmd
	arg1            string
	arg2            int
}

func New(f io.Reader) *Parser {
	s := bufio.NewScanner(f)
	return &Parser{scanner: s, hasMoreCommands: true}
}

func (p *Parser) HasMoreCommands() bool {
	return p.hasMoreCommands
}

func (p *Parser) Advance() error {
	if !p.scanner.Scan() {
		err := p.scanner.Err()
		// io.EOF
		if err == nil {
			p.hasMoreCommands = false
		} else {
			log.Fatal("Parser.Advance: %w", err)
		}
	}
	return p.parse(p.scanner.Bytes())
}

func (p *Parser) CommandType() Cmd {
	return p.currentCmd
}

func (p *Parser) Arg1() string {
	return p.arg1
}

func (p *Parser) Arg2() int {
	return p.arg2
}

var regexpCmd = regexp.MustCompile(`^(?P<cmd>[a-z\-]+)\s*(?P<arg1>[0-9A-Za-z_:\.]+)*\s*(?P<arg2>[0-9]+)*`)

func (p *Parser) parse(row []byte) error {
	p.arg1 = ""
	p.arg2 = 0

	b := bytes.TrimSpace(row)

	// skip empty row
	if len(b) == 0 {
		p.currentCmd = EMPTY
		return nil
	}

	// skip comment
	if bytes.HasPrefix(b, []byte("//")) {
		p.currentCmd = COMMENT
		return nil
	}

	// parse
	matches := regexpCmd.FindSubmatch(b)
	if len(matches) == 0 {
		return fmt.Errorf("invalid format: %s", string(b))
	}

	cmd := matches[regexpCmd.SubexpIndex("cmd")]
	arg1 := matches[regexpCmd.SubexpIndex("arg1")]
	arg2 := matches[regexpCmd.SubexpIndex("arg2")]

	switch string(cmd) {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		p.currentCmd = C_ARITHMETRIC
		p.arg1 = string(cmd)
	case "return":
		p.currentCmd = C_RETURN
	case "label":
		p.currentCmd = C_LABEL
		p.arg1 = string(arg1)
	case "goto":
		p.currentCmd = C_GOTO
		p.arg1 = string(arg1)
	case "if-goto":
		p.currentCmd = C_IF
		p.arg1 = string(arg1)
	// ignore the error because arg2 is ensured that be a numeric by regexp
	case "function":
		p.currentCmd = C_FUNCTION
		p.arg1 = string(arg1)
		p.arg2, _ = strconv.Atoi(string(arg2))
	case "call":
		p.currentCmd = C_CALL
		p.arg1 = string(arg1)
		p.arg2, _ = strconv.Atoi(string(arg2))
	case "push":
		p.currentCmd = C_PUSH
		p.arg1 = string(arg1)
		p.arg2, _ = strconv.Atoi(string(arg2))
	case "pop":
		p.currentCmd = C_POP
		p.arg1 = string(arg1)
		p.arg2, _ = strconv.Atoi(string(arg2))
	default:
		return fmt.Errorf("unknown command: %s %s %s", string(cmd), string(arg1), string(arg2))
	}

	return nil
}
