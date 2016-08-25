package cmd

import (
	"fmt"
	"log"

	"github.com/rybit/config_example/conf"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use: "example",
	Run: run,
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd.PersistentFlags().StringP("config", "c", "", "the config file to use")
	rootCmd.Flags().IntP("port", "p", 0, "the port to use")

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	config, err := conf.LoadConfig(cmd)
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	_ = config
	fmt.Println("--- here ---")
}
