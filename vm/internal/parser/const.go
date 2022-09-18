package parser

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

type Segment string

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
