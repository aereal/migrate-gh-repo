package main

import (
	"log"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Printf("! %v", err)
		os.Exit(1)
	}
}

func run(argv []string) error {
	return nil
}
