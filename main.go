package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	outDir := os.Args[2]
	outDir, err := filepath.Abs(outDir)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(outDir)
	if err != nil {
		err = os.Mkdir(outDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	lib := os.Args[1]

	Gen(lib, outDir, os.Args[3:])
}
