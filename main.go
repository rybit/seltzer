package main

import (
	"log"

	"github.com/rybit/config_example/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
