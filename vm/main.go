package main

import (
	"fmt"
	"log"
	"os"

	"github.com/uu64/nand2tetris/vm/cmd"
)

func usage() {
	fmt.Println("usage: vmc /path/to/file.asm")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	outputFilePath := fmt.Sprintf("%s.asm", os.Args[1][0:len(os.Args[1])-len(".vm")])
	cmd := cmd.New(os.Args[1], outputFilePath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output: %s\n", outputFilePath)
}
