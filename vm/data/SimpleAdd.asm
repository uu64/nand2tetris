// push constant 7
// push constant 8
// add

// init
@256
D=A
@SP
M=D

// push constant 7
@7
D=A
@SP
A=M
M=D
@SP
M=M+1

// push constant 8
@8
D=A
@SP
A=M
M=D
@SP
M=M+1

//add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=D+M
@SP
M=M+1