package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	cmd := cmd.New(os.Args[1], fmt.Sprintf("%s/out.asm", filepath.Dir(ex)))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
