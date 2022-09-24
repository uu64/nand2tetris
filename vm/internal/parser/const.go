package parser

type CmdType int

const (
	COMMENT CmdType = iota
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

const (
	CMD_ADD    = "add"
	CMD_SUB    = "sub"
	CMD_NEG    = "neg"
	CMD_EQ     = "eq"
	CMD_GT     = "gt"
	CMD_LT     = "lt"
	CMD_AND    = "and"
	CMD_OR     = "or"
	CMD_NOT    = "not"
	CMD_RETURN = "return"
	CMD_LABEL  = "label"
	CMD_GOTO   = "goto"
	CMD_IF     = "if-goto"
	CMD_FUNC   = "function"
	CMD_CALL   = "call"
	CMD_PUSH   = "push"
	CMD_POP    = "pop"
)

const (
	SEG_ARG    = "argument"
	SEG_LOCAL  = "local"
	SEG_STATIC = "static"
	SEG_CONST  = "constant"
	SEG_THIS   = "this"
	SEG_THAT   = "that"
	SEG_PTR    = "pointer"
	SEG_TEMP   = "temp"
)
