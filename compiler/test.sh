#!/bin/sh -eu

check_tokenizer() {
    echo "../projects/10/${1}/${2}.jack"
    test -e "../projects/10/${1}/${2}T.xml" && rm "../projects/10/${1}/${2}T.xml"
    ./JackTokenizer "../projects/10/${1}/${2}.jack"
    ../tools/TextComparer.sh "./cmd/tokenizer/data/${1}/${2}T.xml" "../projects/10/${1}/${2}T.xml"
}

test -e ./JackTokenizer && rm ./JackTokenizer
go build -o ./JackTokenizer ./cmd/tokenizer

check_tokenizer "ArrayTest" "Main"

check_tokenizer "ExpressionLessSquare" "Main"
check_tokenizer "ExpressionLessSquare" "Square"
check_tokenizer "ExpressionLessSquare" "SquareGame"

check_tokenizer "Square" "Main"
check_tokenizer "Square" "Square"
check_tokenizer "Square" "SquareGame"


echo "finish."