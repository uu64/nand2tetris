package compile

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/token"
)

type Compiler struct {
	tokenizer *token.Tokenizer
}

func New(t *token.Tokenizer) (*Compiler, error) {
	if err := t.Advance(); err != nil {
		return nil, fmt.Errorf("CompileClass: %w", err)
	}

	return &Compiler{t}, nil
}
