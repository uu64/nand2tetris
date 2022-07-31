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
	arg1            []byte
	arg2            []byte
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

func (p *Parser) Arg1() string {
	return string(p.arg1)
}

func (p *Parser) Arg2() (int, error) {
	return strconv.Atoi(string(p.arg2))
}

var regexpPush = regexp.MustCompile(`^push\s+(?P<arg1>argument|local|static|constant|this|that|pointer|temp)\s+(?P<arg2>[0-9]+)`)
var regexpPop = regexp.MustCompile(`^pop\s+(?P<arg1>argument|local|static|constant|this|that|pointer|temp)\s+(?P<arg2>[0-9]+)`)

var regexpLabel = regexp.MustCompile(`^label\s+(?P<arg1>[0-9A-Za-z_:\.]+)`)
var regexpGoto = regexp.MustCompile(`^goto\s+(?P<arg1>[0-9A-Za-z_:\.]+)`)
var regexpIf = regexp.MustCompile(`^if-goto\s+(?P<arg1>[0-9A-Za-z_:\.]+)`)

var regexpFunc = regexp.MustCompile(`^function\s+(?P<arg1>[0-9A-Za-z_:\.]+)\s+(?P<arg2>[0-9]+)`)
var regexpCall = regexp.MustCompile(`^call\s+(?P<arg1>[0-9A-Za-z_:\.]+)\s+(?P<arg2>[0-9]+)`)

func (p *Parser) parse(row []byte) error {
	p.arg1 = nil
	p.arg2 = nil

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

	switch {
	// TODO: arithmetric
	case bytes.HasPrefix(b, []byte("push")):
		matches := regexpPush.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_PUSH
		p.arg1 = matches[regexpPush.SubexpIndex("arg1")]
		p.arg2 = matches[regexpPush.SubexpIndex("arg2")]
	case bytes.HasPrefix(b, []byte("pop")):
		matches := regexpPop.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_POP
		p.arg1 = matches[regexpPop.SubexpIndex("arg1")]
		p.arg2 = matches[regexpPop.SubexpIndex("arg2")]
	case bytes.HasPrefix(b, []byte("label")):
		matches := regexpLabel.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_LABEL
		p.arg1 = matches[regexpLabel.SubexpIndex("arg1")]
	case bytes.HasPrefix(b, []byte("goto")):
		matches := regexpGoto.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_GOTO
		p.arg1 = matches[regexpGoto.SubexpIndex("arg1")]
	case bytes.HasPrefix(b, []byte("if-goto")):
		matches := regexpIf.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_IF
		p.arg1 = matches[regexpIf.SubexpIndex("arg1")]
	case bytes.HasPrefix(b, []byte("function")):
		matches := regexpFunc.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_FUNCTION
		p.arg1 = matches[regexpCall.SubexpIndex("arg1")]
		p.arg2 = matches[regexpCall.SubexpIndex("arg2")]
	case bytes.HasPrefix(b, []byte("call")):
		matches := regexpCall.FindSubmatch(b)
		if len(matches) == 0 {
			return fmt.Errorf("invalid format: %s", string(b))
		}

		p.currentCmd = C_CALL
		p.arg1 = matches[regexpCall.SubexpIndex("arg1")]
		p.arg2 = matches[regexpCall.SubexpIndex("arg2")]
	case bytes.Equal(b, []byte("return")):
		p.currentCmd = C_RETURN
		p.arg1 = nil
		p.arg2 = nil
	default:
		return fmt.Errorf("unknown command: %s", string(b))
	}
	return nil
}
