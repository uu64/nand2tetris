#!/bin/sh -eu

test -e ./vm && rm ./vm

test() {
    echo "${1}/${2}/${3}"
    ./vm "../projects/${1}/${2}/${3}/${3}.vm"
    ../tools/CPUEmulator.sh "../projects/${1}/${2}/${3}/${3}.tst"
    echo 
}

go build -o ./vm

test "07" "StackArithmetic" "SimpleAdd"
test "07" "StackArithmetic" "StackTest"

test "07" "MemoryAccess" "BasicTest"
test "07" "MemoryAccess" "PointerTest"
test "07" "MemoryAccess" "StaticTest"

test "08" "ProgramFlow" "BasicLoop"
test "08" "ProgramFlow" "FibonacciSeries"

test "08" "FunctionCalls" "SimpleFunction"