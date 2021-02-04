package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/common"
)

var participantsFilepath string
var instructionsFilepath string
var participants []*common.Participant
var instructionsTMPL *template.Template

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
	rootCmd.PersistentFlags().StringVarP(&participantsFilepath, "participants", "p", "", "input file containing participants")
	rootCmd.PersistentFlags().StringVarP(&instructionsFilepath, "instructions", "i", "", "input file containing instructions template")
}

// Execute executes a CLI command - boilerplate for cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
