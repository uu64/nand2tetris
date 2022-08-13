package code

// add two values.
const Add = `@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1
`

// push the value to a constant segment.
const PushConstant = `@%d
D=A
@SP
A=M
M=D
@SP
M=M+1
`

// pop from a constant segment.
// M is set to the return value.
const PopConstant = `@SP
AM=M-1
`
