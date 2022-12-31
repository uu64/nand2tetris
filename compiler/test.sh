#!/bin/bash -eu

check_tokenizer() {
    echo "../projects/10/${1}/${2}.jack"
    test -e "../projects/10/${1}/${2}T.xml" && rm "../projects/10/${1}/${2}T.xml"
    ./JackTokenizer "../projects/10/${1}/${2}.jack"
    diff -uw "./cmd/tokenizer/data/${1}/${2}T.xml" "../projects/10/${1}/${2}T.xml"
}

check_compiler() {
    echo "../projects/10/${1}/${2}.jack"
    test -e "../projects/10/${1}/${2}.xml" && rm "../projects/10/${1}/${2}.xml"
    ./JackCompiler "../projects/10/${1}/${2}.jack"
    diff -uw  <(gsed -z -E "s#<([a-zA-Z]+)>\r?\n\s+</([a-zA-Z]+)>#<\1></\2>#g" "./cmd/compiler/data/${1}/${2}.xml") "../projects/10/${1}/${2}.xml"
}

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
test JackTokenizer
####################

EOF

check_tokenizer "ArrayTest" "Main"

check_tokenizer "ExpressionLessSquare" "Main"
check_tokenizer "ExpressionLessSquare" "Square"
check_tokenizer "ExpressionLessSquare" "SquareGame"

check_tokenizer "Square" "Main"
check_tokenizer "Square" "Square"
check_tokenizer "Square" "SquareGame"


cat<<EOF

####################
test JackCompiler
####################

EOF

check_compiler "ExpressionLessSquare" "Main"
check_compiler "ExpressionLessSquare" "Square"
check_compiler "ExpressionLessSquare" "SquareGame"

echo "finish."