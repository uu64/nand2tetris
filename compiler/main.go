package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/uu64/nand2tetris/compiler/cmd"
)

func usage() {
	fmt.Println("usage: jackc input")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		return
	}

	// TODO: directoryかファイル単体を渡す
	cmd := cmd.New(args[0])
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("finish.\n")
}
