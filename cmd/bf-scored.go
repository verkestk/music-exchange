package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/bfscored"
	"github.com/verkestk/music-exchange/common"
	"github.com/verkestk/music-exchange/email"
)

func init() {
	pairRootCmd.AddCommand(bfScoredCmd)
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
		participants, err = common.GetParticipantsFromJSONFile(participantsFilepath, true)
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
		pairs, err := bfscored.DoExchange(participants, instructionsTMPL)
		if err != nil {
			return err
		}

		if emailInstructions {
			sender := email.GetSMTPSender(smtpHostEnvVar, smtpPortEnvVar, smtpUsernameEnvVar, smtpPasswordEnvVar)
			err = common.EmailInstructions(pairs, instructionsTMPL, emailSubject, emailTestRecipient, sender)
		} else {
			err = common.WriteInstructions(pairs, instructionsTMPL, instructionsFileExtension)
		}

		if err != nil {
			return err
		}

		return common.UpdateParticipantsJSON(participantsFilepath, pairs)
	},
}
