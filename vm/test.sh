#!/bin/sh -eu

test() {
    echo "${1}/${2}"
    go run ./main.go "../projects/07/${1}/${2}/${2}.vm"
    ../tools/CPUEmulator.sh "../projects/07/${1}/${2}/${2}.tst"
    echo 
}

test "StackArithmetic" "SimpleAdd"
test "StackArithmetic" "StackTest"

test "MemoryAccess" "BasicTest"
