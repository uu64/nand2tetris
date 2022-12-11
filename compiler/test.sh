#!/bin/sh -eu

check() {
    echo "../projects/10/${1}/${2}.jack"
    test -e "../projects/10/${1}/${2}T.xml" && rm "../projects/10/${1}/${2}T.xml"
    ./jackc "../projects/10/${1}/${2}.jack"
    diff -u --strip-trailing-cr "./cmd/tokenizer/data/${1}/${2}T.xml" "../projects/10/${1}/${2}T.xml"
}

test -e ./jackc && rm ./jackc
go build -o ./jackc ./cmd/tokenizer

check "ArrayTest" "Main"

check "ExpressionLessSquare" "Main"
check "ExpressionLessSquare" "Square"
check "ExpressionLessSquare" "SquareGame"

check "Square" "Main"
check "Square" "Square"
check "Square" "SquareGame"


echo "finish."