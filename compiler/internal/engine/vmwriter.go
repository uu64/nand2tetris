package engine

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
	"github.com/uu64/nand2tetris/compiler/internal/vmwriter"
)

func (c *Compiler) writeFunction(kwd *tokenizer.Keyword, name *tokenizer.Identifier, nLocals int) error {
	switch kwd.Label {
	case "constructor":
		return c.codewriter.WriteFunction(fmt.Sprintf("%s.%s", c.symtab.ClassName, name.Label), nLocals)
	case "function":
		return c.codewriter.WriteFunction(fmt.Sprintf("%s.%s", c.symtab.ClassName, name.Label), nLocals)
	case "method":
		return c.codewriter.WriteFunction(name.Label, nLocals)
	default:
		return fmt.Errorf("writeFunction: unexpected keyword %s", kwd.Label)
	}
}

func (c *Compiler) writeOp(op Op) error {
	tk := tokenizer.Symbol(op)
	switch tk.Val() {
	case '+':
		return c.codewriter.WriteArithmetic(vmwriter.Add)
	case '-':
		return c.codewriter.WriteArithmetic(vmwriter.Sub)
	case '*':
		return c.codewriter.WriteCall("Math.multiply", 2)
	case '/':
		return c.codewriter.WriteCall("Math.divide", 2)
	default:
		return fmt.Errorf("writeOp: undefined op %s", string(tk.Val()))
	}
}

func (c *Compiler) writeExpression(exp *Expression) error {
	for _, e := range exp.Tokens {
		fmt.Println(e.ElementType().String())
	}
	fmt.Println()
	return nil
}

func (c *Compiler) writeSubroutine(el *SubroutineDec) error {
	for _, tk := range el.Tokens {
		switch tk.ElementType() {
		case tokenizer.ElToken:
			tk := tk.(tokenizer.Token)
			fmt.Println(tk.TokenType().String())
		default:
			fmt.Println(tk.ElementType())
		}
	}
	return nil
}
