// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

(MAIN)
    @i
    M=0
    @offset
    M=0

    // check key pressed or not
    @24576
    D=M
    @key_pressed
    M=D
(VLOOP_START)
    @i
    D=M
    @256
    D=D-A
    @VLOOP_END
    D;JGE

    @j
    M=0
(HLOOP_START)
    @j
    D=M
    @32
    D=D-A
    @HLOOP_END
    D;JGE

    // clear or fill
    @key_pressed
    D=M
    @CLEAR
    D;JEQ
    @FILL
    0;JMP
(CLEAR)
    @offset
    D=M
    @SCREEN
    A=A+D
    // reset
    M=0
    @CONTINUE
    0;JMP
(FILL)
    @offset
    D=M
    @SCREEN
    A=A+D
    // fill
    M=-1
    @CONTINUE
    0;JMP

(CONTINUE)
    // add offset
    @offset
    M=M+1
    // add loop variable 
    @j
    M=M+1
    @HLOOP_START
    0;JMP
(HLOOP_END)
    // add loop variable 
    @i
    M=M+1

    // go to VLOOP_START
    @VLOOP_START
    D;JMP
(VLOOP_END)
    // go to MAIN
    @MAIN
    0;JMP
