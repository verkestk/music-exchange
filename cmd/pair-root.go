package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/common"
)

var (
	participantsFilepath      string
	instructionsFilepath      string
	participants              []*common.Participant
	instructionsTMPL          *template.Template
	emailInstructions         bool
	instructionsFileExtension string
	updateParticipantsFile    bool
	emailSubject              string
	emailTestRecipient        string
	smtpHostEnvVar            string
	smtpPortEnvVar            string
	smtpUsernameEnvVar        string
	smtpPasswordEnvVar        string
)

func init() {
	rootCmd.AddCommand(pairRootCmd)
	pairRootCmd.PersistentFlags().StringVarP(&participantsFilepath, "participants", "p", "", "input file containing participants")
	pairRootCmd.PersistentFlags().StringVarP(&instructionsFilepath, "instructions", "i", "", "input file containing instructions template")
	pairRootCmd.PersistentFlags().BoolVarP(&emailInstructions, "email", "e", false, "set in order to email the instructions rather than generate local files")
	pairRootCmd.PersistentFlags().StringVarP(&instructionsFileExtension, "extension", "x", "md", "the extension to use when writing instructions to local files")
	pairRootCmd.PersistentFlags().BoolVarP(&updateParticipantsFile, "update", "u", false, "set in order to update the JSON input file with the latest pairings")
	pairRootCmd.PersistentFlags().StringVarP(&emailSubject, "subject", "b", "Music Exchange Assignment", "set this to change the subject the emails are sent with")
	pairRootCmd.PersistentFlags().StringVarP(&emailTestRecipient, "recipient", "r", "", "set this to always send the instructions to this specific email address - great for testing")
	pairRootCmd.PersistentFlags().StringVarP(&smtpHostEnvVar, "smtp-host", "", "MUSIC_EXCHANGE_SMTP_HOST", "change the environment variable this program references to get the SMPT host")
	pairRootCmd.PersistentFlags().StringVarP(&smtpPortEnvVar, "smtp-port", "", "MUSIC_EXCHANGE_SMTP_PORT", "change the environment variable this program references to get the SMPT port")
	pairRootCmd.PersistentFlags().StringVarP(&smtpUsernameEnvVar, "smtp-username", "", "MUSIC_EXCHANGE_SMTP_USERNAME", "change the environment variable this program references to get the SMPT username")
	pairRootCmd.PersistentFlags().StringVarP(&smtpPasswordEnvVar, "smtp-password", "", "MUSIC_EXCHANGE_SMTP_PASSWORD", "change the environment variable this program references to get the SMPT password")
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
