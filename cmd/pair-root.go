package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/common"
)

var participantsFilepath string
var instructionsFilepath string
var participants []*common.Participant
var instructionsTMPL *template.Template

func init() {
	rootCmd.AddCommand(pairRootCmd)
	pairRootCmd.PersistentFlags().StringVarP(&participantsFilepath, "participants", "p", "", "input file containing participants")
	pairRootCmd.PersistentFlags().StringVarP(&instructionsFilepath, "instructions", "i", "", "input file containing instructions template")
	pairRootCmd.MarkFlagRequired("participants")
	pairRootCmd.MarkFlagRequired("instructions")
}

var pairRootCmd = &cobra.Command{
	Use:   "pair",
	Short: "pairs participants based on selected algorithm.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You must select a subcommand!\n Run the `help` command for instructions.")
	},
}
