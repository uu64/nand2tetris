#!/bin/sh -eu

test -e ./vmc && rm ./vmc

check() {
    echo "${1}/${2}/${3}"
    test -e "../projects/${1}/${2}/${3}/${3}.asm" && rm "../projects/${1}/${2}/${3}/${3}.asm"
    find "../projects/${1}/${2}/${3}" -name "*vm" -exec ./vmc -o "../projects/${1}/${2}/${3}/${3}.asm" {} +
    ../tools/CPUEmulator.sh "../projects/${1}/${2}/${3}/${3}.tst"
    echo 
}

go build -o ./vmc

check "07" "StackArithmetic" "SimpleAdd"
check "07" "StackArithmetic" "StackTest"

check "07" "MemoryAccess" "BasicTest"
check "07" "MemoryAccess" "PointerTest"
check "07" "MemoryAccess" "StaticTest"

check "08" "ProgramFlow" "BasicLoop"
check "08" "ProgramFlow" "FibonacciSeries"

check "08" "FunctionCalls" "SimpleFunction"