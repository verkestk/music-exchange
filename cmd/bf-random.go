package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/bfrandom"
	"github.com/verkestk/music-exchange/common"
)

var avoid int

func init() {
	rootCmd.AddCommand(bfRandomCmd)
	bfRandomCmd.Flags().IntVarP(&avoid, "avoid", "a", 0, "how many past recipients to avoid")
	bfRandomCmd.MarkFlagRequired("filepath")
}

var bfRandomCmd = &cobra.Command{
	Use:   "bfRandom",
	Short: "Shuffle the particpants until all conditions are satisfied",
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
		return bfrandom.DoExchange(participants, avoid)
	},
}
