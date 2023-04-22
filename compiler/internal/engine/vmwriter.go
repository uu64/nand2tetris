package engine

import (
	"fmt"

	"github.com/uu64/nand2tetris/compiler/internal/tokenizer"
	"github.com/uu64/nand2tetris/compiler/internal/vmwriter"
)

func (c *Compiler) writeFunction(kwd, name string, nLocals int) error {
	switch kwd {
	case "constructor":
		return c.codewriter.WriteFunction(fmt.Sprintf("%s.%s", c.ctx.ClassName, name), nLocals)
	case "function":
		return c.codewriter.WriteFunction(fmt.Sprintf("%s.%s", c.ctx.ClassName, name), nLocals)
	case "method":
		return c.codewriter.WriteFunction(name, nLocals)
	default:
		return fmt.Errorf("writeFunction: unexpected keyword %s", kwd)
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
	case '&':
		return c.codewriter.WriteArithmetic(vmwriter.And)
	case '|':
		return c.codewriter.WriteArithmetic(vmwriter.Or)
	case '<':
		return c.codewriter.WriteArithmetic(vmwriter.LT)
	case '>':
		return c.codewriter.WriteArithmetic(vmwriter.GT)
	case '=':
		return c.codewriter.WriteArithmetic(vmwriter.EQ)
	default:
		panic(fmt.Sprintf("writeOp: undefined op %s", string(tk.Val())))
	}
}

func (c *Compiler) writeUnaryOp(op UnaryOp) error {
	tk := tokenizer.Symbol(op)
	switch tk.Val() {
	case '-':
		return c.codewriter.WriteArithmetic(vmwriter.Neg)
	case '~':
		return c.codewriter.WriteArithmetic(vmwriter.Not)
	default:
		panic(fmt.Sprintf("writeUnaryOp: undefined op %s", string(tk.Val())))
	}
}

func (c *Compiler) writeCall(name string, nArgs int) {
	c.codewriter.WriteCall(name, nArgs)
}

func (c *Compiler) writePushIntConst(n int) {
	c.codewriter.WritePush(vmwriter.Const, n)
}

func (c *Compiler) writeReturn(isVoid bool) {
	if isVoid {
		c.codewriter.WritePush(vmwriter.Const, 0)
	}
	c.codewriter.WriteReturn()
}

func (c *Compiler) discardReturn() {
	c.codewriter.WritePop(vmwriter.Temp, 0)
}
