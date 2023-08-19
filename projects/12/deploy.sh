#!/bin/bash

function deploy() {
    cp "${1}.jack" "${1}Test/"
    # cp ../../tools/OS/* "${1}Test/"
    ../../tools/JackCompiler.sh "${1}Test/"
}

deploy "Array"
deploy "Keyboard"
deploy "Math"
# deploy "Memory"
deploy "Output"
deploy "Screen"
deploy "String"
deploy "Sys"