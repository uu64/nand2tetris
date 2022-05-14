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
@FILL
0;JMP

(CLEAR)
// clear screen
    @i
    M=0
    @offset
    M=0
(CLEAR_VLOOP_START)
    @i
    D=M
    @256
    D=D-A
    @CLEAR_VLOOP_END
    D;JGE

    @j
    M=0
(CLEAR_HLOOP_START)
    @j
    D=M
    @32
    D=D-A
    @CLEAR_HLOOP_END
    D;JGE

    // calculate address
    @offset
    D=M
    @SCREEN
    A=A+D

    // clear
    M=0

    // add offset
    @offset
    M=M+1
    // add loop variable 
    @j
    M=M+1
    @CLEAR_HLOOP_START
    0;JMP
(CLEAR_HLOOP_END)
    // add loop variable 
    @i
    M=M+1

    // go to CLEAR_VLOOP_START
    @CLEAR_VLOOP_START
    D;JMP
(CLEAR_VLOOP_END)
    // go to MAIN
    @MAIN
    0;JMP

(FILL)
// fill screen
    @i
    M=0
    @offset
    M=0
(FILL_VLOOP_START)
    @i
    D=M
    @256
    D=D-A
    @FILL_VLOOP_END
    D;JGE

    @j
    M=0
(FILL_HLOOP_START)
    @j
    D=M
    @32
    D=D-A
    @FILL_HLOOP_END
    D;JGE

    // calculate address
    @offset
    D=M
    @SCREEN
    A=A+D

    // fill
    M=-1

    // add offset
    @offset
    M=M+1
    // add loop variable 
    @j
    M=M+1
    @FILL_HLOOP_START
    0;JMP
(FILL_HLOOP_END)
    // add loop variable 
    @i
    M=M+1

    // go to FILL_VLOOP_START
    @FILL_VLOOP_START
    D;JMP
(FILL_VLOOP_END)
    // go to MAIN
    @MAIN
    0;JMP