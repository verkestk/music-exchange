package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/bfscored"
	"github.com/verkestk/music-exchange/common"
)

func init() {
	rootCmd.AddCommand(bfScoredCmd)
	bfScoredCmd.MarkFlagRequired("participants")
	bfScoredCmd.MarkFlagRequired("instructions")
}

var bfScoredCmd = &cobra.Command{
	Use:   "bfScored",
	Short: "Evaluate all possible pairings and pick one of the best",
	Args: func(cmd *cobra.Command, args []string) error {
		if participantsFilepath == "" {
			return fmt.Errorf("participants required")
		}

		if instructionsFilepath == "" {
			return fmt.Errorf("instructions required")
		}

		var err error
		participants, err = common.GetParticipantsFromFile(participantsFilepath)
		if err != nil {
			return err
		}

		instructionsTMPL, err = template.ParseFiles(instructionsFilepath)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return bfscored.DoExchange(participants, instructionsTMPL)
	},
}
