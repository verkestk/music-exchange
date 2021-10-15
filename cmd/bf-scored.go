package cmd

import (
	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/operation"
)

func init() {
	pairRootCmd.AddCommand(bfScoredCmd)
}

var bfScoredCmd = &cobra.Command{
	Use:   "bfScored",
	Short: "Evaluate all possible pairings and pick one of the best",
	Args: func(cmd *cobra.Command, args []string) error {
		pairConfig.Algorithm = operation.BFScored
		pairConfig.EmailSender = email.GetGmailSender(smtpHostEnvVar, smtpPortEnvVar, smtpUsernameEnvVar, smtpPasswordEnvVar)
		return pairConfig.Prepare()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return operation.DoPair(pairConfig)
	},
}
