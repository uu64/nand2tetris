#!/bin/bash

function deploy() {
    cp "${1}.jack" "${2}/"
    ../../tools/JackCompiler.sh "${2}/"
}

deploy "Array" "ArrayTest"
deploy "Keyboard" "KeyboardTest"
deploy "Math" "MathTest"
deploy "Memory" "MemoryTest"
# deploy "Memory" "MemoryTest/MemoryDiag"
deploy "Output" "OutputTest"
deploy "Screen" "ScreenTest"
deploy "String" "StringTest"
deploy "Sys" "SysTest"