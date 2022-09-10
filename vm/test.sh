#!/bin/sh -eu

go run ./main.go ../projects/07/StackArithmetic/SimpleAdd/SimpleAdd.vm
../tools/CPUEmulator.sh ../projects/07/StackArithmetic/SimpleAdd/SimpleAdd.tst