package cmd

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/bfrandom"
	"github.com/verkestk/music-exchange/common"
)

var avoid int

func init() {
	pairRootCmd.AddCommand(bfRandomCmd)
	bfRandomCmd.Flags().IntVarP(&avoid, "avoid", "a", 0, "how many past recipients to avoid")
}

var bfRandomCmd = &cobra.Command{
	Use:   "bfRandom",
	Short: "Shuffle the particpants until all conditions are satisfied",
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
		pairs, err := bfrandom.DoExchange(participants, instructionsTMPL, avoid)
		if err != nil {
			return err
		}

		if emailInstructions {
			err = common.EmailInstructions(pairs, instructionsTMPL, emailSubject, emailTestRecipient, smtpHostEnvVar, smtpPortEnvVar, smtpUsernameEnvVar, smtpPasswordEnvVar)
		} else {
			err = common.WriteInstructions(pairs, instructionsTMPL, instructionsFileExtension)
		}

		if err != nil {
			return err
		}

		return common.UpdateParticipantsJSON(participantsFilepath, pairs)
	},
}
