package cmd

import (
	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/operation"
	"github.com/verkestk/music-exchange/src/email"
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
		pairConfig.EmailSender = email.GetGmailSender(smtpHostEnvVar, smtpPortEnvVar, smtpUsernameEnvVar, smtpPasswordEnvVar)
		return pairConfig.Prepare()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := operation.DoPair(pairConfig)
		return err
	},
}
