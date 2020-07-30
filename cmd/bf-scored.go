package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/bfscored"
	"github.com/verkestk/music-exchange/common"
)

func init() {
	rootCmd.AddCommand(bfScoredCmd)
	bfScoredCmd.MarkFlagRequired("filepath")
}

var bfScoredCmd = &cobra.Command{
	Use:   "bfScored",
	Short: "Evaluate all possible pairings and pick one of the best",
	Args: func(cmd *cobra.Command, args []string) error {
		if filepath == "" {
			return fmt.Errorf("filepath required")
		}

		var err error
		participants, err = common.GetParticipantsFromFile(filepath)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return bfscored.DoExchange(participants)
	},
}
