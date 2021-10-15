package cmd

import (
	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/operation"
)

func init() {
	pairRootCmd.AddCommand(bfRandomCmd)
	bfRandomCmd.Flags().IntVarP(&pairConfig.Avoid, "avoid", "a", 0, "how many past recipients to avoid")
}

var bfRandomCmd = &cobra.Command{
	Use:   "bfRandom",
	Short: "Shuffle the particpants until all conditions are satisfied",
	Args: func(cmd *cobra.Command, args []string) error {
		pairConfig.Algorithm = operation.BFRandom
		return pairConfig.Prepare()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return operation.DoPair(pairConfig)
	},
}
