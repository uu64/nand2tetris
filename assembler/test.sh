#!/bin/sh -eu

test() {
    [ -e "${1}/../projects/06/${2}/${3}.hack" ] && rm "${1}/../projects/06/${2}/${3}.hack"
    [ -e "${1}/.hack" ] && rm "${1}/.hack"

    echo "${2}/${3}.asm"

    go run "${1}/main.go" "${1}/../projects/06/${2}/${3}.asm"
    "${1}/../tools/Assembler.sh" "${1}/../projects/06/${2}/${3}.asm"

    diff -u "${1}/.hack" "${1}/../projects/06/${2}/${3}.hack"

    echo "no diff"
    echo
}

test "$(dirname "${0}")" "add" "Add"
test "$(dirname "${0}")" "max" "Max"
test "$(dirname "${0}")" "max" "MaxL"
test "$(dirname "${0}")" "pong" "Pong"
test "$(dirname "${0}")" "pong" "PongL"