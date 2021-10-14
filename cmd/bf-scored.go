package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/bfscored"
	"github.com/verkestk/music-exchange/common"
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
			err = common.EmailInstructions(pairs, instructionsTMPL, emailSubject, emailTestRecipient, smtpHostEnvVar, smtpPortEnvVar, smtpUsernameEnvVar, smtpPasswordEnvVar)
		} else {
			err = common.WriteInstructions(pairs, instructionsTMPL, instructionsFileExtension)
		}

		return err
	},
}
