package codewriter

import "fmt"

var opCounter = 0

// eq compares the two values
// and sets -1 if they are equal and 0 otherwise in the stack
const eqTmpl = `@SP
AM=M-1
D=M
@SP
AM=M-1
D=M-D
@EQ%[1]d_T
D;JEQ
@EQ%[1]d_F
0;JMP
(EQ%[1]d_T)
D=-1
@EQ%[1]d_END
0;JMP
(EQ%[1]d_F)
D=0
@EQ%[1]d_END
0;JMP
(EQ%[1]d_END)
@SP
A=M
M=D
@SP
M=M+1
`

func eq() string {
	opCounter += 1
	return fmt.Sprintf(eqTmpl, opCounter)
}

// gt compares the two values
// and sets -1 if x > y and 0 otherwise in the stack
const gtTmpl = `@SP
AM=M-1
D=M
@SP
AM=M-1
D=M-D
@GT%[1]d_T
D;JGT
@GT%[1]d_F
0;JMP
(GT%[1]d_T)
D=-1
@GT%[1]d_END
0;JMP
(GT%[1]d_F)
D=0
@GT%[1]d_END
0;JMP
(GT%[1]d_END)
@SP
A=M
M=D
@SP
M=M+1
`

func gt() string {
	opCounter += 1
	return fmt.Sprintf(gtTmpl, opCounter)
}

// lt compares the two values
// and sets -1 if x < y and 0 otherwise in the stack
const ltTmpl = `@SP
AM=M-1
D=M
@SP
AM=M-1
D=M-D
@LT%[1]d_T
D;JLT
@LT%[1]d_F
0;JMP
(LT%[1]d_T)
D=-1
@LT%[1]d_END
0;JMP
(LT%[1]d_F)
D=0
@LT%[1]d_END
0;JMP
(LT%[1]d_END)
@SP
A=M
M=D
@SP
M=M+1
`

func lt() string {
	opCounter += 1
	return fmt.Sprintf(ltTmpl, opCounter)
}
