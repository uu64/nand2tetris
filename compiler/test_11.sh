#!/bin/bash -eu

check_compiler() {
    echo "../projects/10/${1}/${2}.jack"
    test -e "../projects/10/${1}/${2}.xml" && rm "../projects/10/${1}/${2}.xml"
    ./JackCompiler "../projects/10/${1}/${2}.jack"
    # replace: <tag>{new line}{indent}</tag> -> <tag></tag>
    diff -uw  <(gsed -z -E "s#<([a-zA-Z]+)>\r?\n\s+</([a-zA-Z]+)>#<\1></\2>#g" "./cmd/compiler/data/${1}/${2}.xml") "../projects/10/${1}/${2}.xml"
    echo "pass"
}

go version

cat<<EOF

####################
compile module
####################

EOF

test -e ./JackTokenizer && rm ./JackTokenizer
go build -o ./JackTokenizer ./cmd/tokenizer

test -e ./JackCompiler && rm ./JackCompiler
go build -o ./JackCompiler ./cmd/compiler

echo "compile success"

cat<<EOF

####################
test JackCompiler
####################

EOF

check_compiler "ExpressionLessSquare" "Main"
check_compiler "ExpressionLessSquare" "Square"
check_compiler "ExpressionLessSquare" "SquareGame"

check_compiler "ArrayTest" "Main"

check_compiler "Square" "Main"
check_compiler "Square" "Square"
check_compiler "Square" "SquareGame"

echo "finish."