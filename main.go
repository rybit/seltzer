package main

import (
	"log"

	"github.com/rybit/seltzer/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
