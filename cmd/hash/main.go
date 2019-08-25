package main

import (
	"encoder-backend/pkg/file"
	"flag"
	"fmt"
)

func main() {

	src := flag.String("file", "", "the absolute path to the file to be hashed")

	flag.Parse()

	fmt.Printf("Checking: %s\n", *src)

	hash, err := file.Checksum(*src)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Hash: %s\n", hash)
}
