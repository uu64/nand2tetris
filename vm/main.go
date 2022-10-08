package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/uu64/nand2tetris/vm/cmd"
)

var output = flag.String("o", "out.asm", "output file")

func usage() {
	fmt.Println("usage: vmc [-o output] input [input ...]")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	flag.Parse()

	cmd := cmd.New(flag.Args(), *output)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output: %s\n", *output)
}
