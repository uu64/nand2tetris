#!/bin/sh -eu

test -e ./vmc && rm ./vmc

check() {
    echo "${1}/${2}/${3}"
    test -e "../projects/${1}/${2}/${3}/${3}.asm" && rm "../projects/${1}/${2}/${3}/${3}.asm"
    find "../projects/${1}/${2}/${3}" -name "*vm" -exec ./vmc -o "../projects/${1}/${2}/${3}/${3}.asm" {} +
    ../tools/CPUEmulator.sh "../projects/${1}/${2}/${3}/${3}.tst"
    echo 
}

check_with_noboot() {
    echo "${1}/${2}/${3}"
    test -e "../projects/${1}/${2}/${3}/${3}.asm" && rm "../projects/${1}/${2}/${3}/${3}.asm"
    find "../projects/${1}/${2}/${3}" -name "*vm" -exec ./vmc -noboot -o "../projects/${1}/${2}/${3}/${3}.asm" {} +
    ../tools/CPUEmulator.sh "../projects/${1}/${2}/${3}/${3}.tst"
    echo 
}

go build -o ./vmc

check_with_noboot "07" "StackArithmetic" "SimpleAdd"
check_with_noboot "07" "StackArithmetic" "StackTest"

check_with_noboot "07" "MemoryAccess" "BasicTest"
check_with_noboot "07" "MemoryAccess" "PointerTest"
check_with_noboot "07" "MemoryAccess" "StaticTest"

check_with_noboot "08" "ProgramFlow" "BasicLoop"
check_with_noboot "08" "ProgramFlow" "FibonacciSeries"

check_with_noboot "08" "FunctionCalls" "SimpleFunction"

check "08" "FunctionCalls" "FibonacciElement"
check "08" "FunctionCalls" "StaticsTest"
check "08" "FunctionCalls" "NestedCall"