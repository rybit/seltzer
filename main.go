package main

import (
	"log"

	"github.com/netlify/netlify-subscriptions/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
