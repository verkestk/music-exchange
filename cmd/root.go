package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/common"
)

var filepath string
var participants []*common.Participant

var rootCmd = &cobra.Command{
	Use:   "musex",
	Short: "Musex provides a secret-santa-style music exchange.",
	Long: `A CLI for pairing music lovers together for a secret-santa-style
                music exchange. Complete documentation at
                https://github.com/verkestk/music-exchange`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to musix!\n Run the `help` command for instructions.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filepath, "filepath", "f", "", "input file containing participants")
}

// Execute executes a CLI command - boilerplate for cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
