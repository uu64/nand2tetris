#!/bin/sh -eu

test -e ./vm && rm ./vm

test() {
    echo "${1}/${2}"
    ./vm "../projects/07/${1}/${2}/${2}.vm"
    ../tools/CPUEmulator.sh "../projects/07/${1}/${2}/${2}.tst"
    echo 
}

go build -o ./vm

test "StackArithmetic" "SimpleAdd"
test "StackArithmetic" "StackTest"

test "MemoryAccess" "BasicTest"
test "MemoryAccess" "PointerTest"
test "MemoryAccess" "StaticTest"
