package main

import (
	"fmt"
	"log"
	"os"

	"github.com/uu64/nand2tetris/assembler/cmd"
)

func usage() {
	fmt.Println("usage: asmc /path/to/file.asm")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	cmd := cmd.New(os.Args[1])
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
