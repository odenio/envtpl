package main

import (
	"log"

	"github.com/odenio/envtpl/v2"
	"github.com/spf13/cobra"
)

// options
var (
	missingKey = "default"
	output     string
	version    bool
)

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var RootCmd = &cobra.Command{
	Use:   "envtpl",
	Short: "Render go templates from environment variables",
	Long:  `Render go templates from environment variables.`,
	Run:   envtpl.Template,
}

func main() {
	RootCmd.Flags().StringVarP(&missingKey, "missingkey", "m", missingKey, "Strategy for dealing with missing keys: [default|zero|error]")
	RootCmd.Flags().StringVarP(&output, "output", "o", output, "The rendered output file")
	RootCmd.Flags().BoolVarP(&version, "version", "v", false, "Prints version")

	err := RootCmd.Execute()
	die(err)
}
