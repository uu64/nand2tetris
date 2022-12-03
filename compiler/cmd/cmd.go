package cmd

import (
	"os"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Cmd struct {
	source string
}

func New(source string) *Cmd {
	return &Cmd{source}
}

func (cmd *Cmd) Run() (err error) {
	f, err := os.Open(cmd.source)
	if err != nil {
		return err
	}
	defer f.Close()

	t := token.New(f)
	for t.HasMoreTokens() {
		t.Advance()
	}

	return nil
}
