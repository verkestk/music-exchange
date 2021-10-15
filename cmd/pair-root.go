package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/operation"
)

var (
	pairConfig = &operation.PairConfig{}
)

func init() {
	rootCmd.AddCommand(pairRootCmd)
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.ParticipantsFilepath, "participants", "p", "", "input file containing participants")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.InstructionsFilepath, "instructions", "i", "", "input file containing instructions template")
	pairRootCmd.PersistentFlags().BoolVarP(&pairConfig.EmailInstructions, "email", "e", false, "set in order to email the instructions rather than generate local files")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.InstructionsFileExtension, "extension", "x", "md", "the extension to use when writing instructions to local files")
	pairRootCmd.PersistentFlags().BoolVarP(&pairConfig.UpdateParticipantsFile, "update", "u", false, "set in order to update the JSON input file with the latest pairings")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.EmailSubject, "subject", "b", "Music Exchange Assignment", "set this to change the subject the emails are sent with")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.EmailTestRecipient, "recipient", "r", "", "set this to always send the instructions to this specific email address - great for testing")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.SMTPHostEnvVar, "smtp-host", "", "MUSIC_EXCHANGE_SMTP_HOST", "change the environment variable this program references to get the SMPT host")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.SMTPPortEnvVar, "smtp-port", "", "MUSIC_EXCHANGE_SMTP_PORT", "change the environment variable this program references to get the SMPT port")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.SMTPUsernameEnvVar, "smtp-username", "", "MUSIC_EXCHANGE_SMTP_USERNAME", "change the environment variable this program references to get the SMPT username")
	pairRootCmd.PersistentFlags().StringVarP(&pairConfig.SMTPPasswordEnvVar, "smtp-password", "", "MUSIC_EXCHANGE_SMTP_PASSWORD", "change the environment variable this program references to get the SMPT password")
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
