package main

import (
	"log"
	"os"

	"github.com/aereal/migrate-gh-repo/config"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Printf("! %v", err)
		os.Exit(1)
	}
}

func run(argv []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log.Printf("config = %#v", cfg)
	return nil
}
