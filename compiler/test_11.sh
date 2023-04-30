#!/bin/bash -eu

check_compiler() {
    echo "../projects/11/${1}/${2}.jack"
    test -e "./cmd/compiler/data/11/${1}/${2}.vm" && rm "./cmd/compiler/data/11/${1}/${2}.vm"
    ./JackCompiler "./cmd/compiler/data/11/${1}/${2}.jack"
    diff -uw  "../projects/11/${1}/${2}.vm" "./cmd/compiler/data/11/${1}/${2}.vm" 
    echo "pass"
}

go version

cat<<EOF

####################
compile module
####################

EOF

test -e ./JackCompiler && rm ./JackCompiler
go build -o ./JackCompiler ./cmd/compiler

echo "compile success"

cat<<EOF

####################
test JackCompiler
####################

EOF

check_compiler "Seven" "Main"

check_compiler "ConvertToBin" "Main"

check_compiler "Square" "Main"
check_compiler "Square" "Square"
check_compiler "Square" "SquareGame"

check_compiler "Average" "Main"

echo "finish."