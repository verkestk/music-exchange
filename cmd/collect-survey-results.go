package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/common"
)

var (
	previousParticipantsFilepath string
	previousParticipants         []*common.Participant
	surveyFilepath               string
	newParticipants              []*common.Participant
	emailAddressColumn           int
	platformsColumn              int
	ignoreColumnsStr             string
	ignoreColumns                []int
	platformsSeparator           string
)

func init() {
	rootCmd.AddCommand(collectSuveyResults)
	collectSuveyResults.Flags().StringVarP(&surveyFilepath, "survey", "s", "", "filepath of CSV survey results")
	collectSuveyResults.Flags().StringVarP(&previousParticipantsFilepath, "previous-participants", "p", "", "filepath of the JSON file used for the last exchange")
	collectSuveyResults.Flags().IntVarP(&emailAddressColumn, "email-address", "e", -1, "index of the column containing the email address")
	collectSuveyResults.Flags().IntVarP(&platformsColumn, "platforms", "l", -1, "index of the column containing the music platforms")
	collectSuveyResults.Flags().StringVarP(&ignoreColumnsStr, "ignore", "i", "", "indexes of columns to ignore")
	collectSuveyResults.Flags().StringVarP(&platformsSeparator, "separator", "a", ";", "separator character between platform choices")
	collectSuveyResults.MarkFlagRequired("survey")
	collectSuveyResults.MarkFlagRequired("email-address")
	collectSuveyResults.MarkFlagRequired("platformsColumn")
}

var collectSuveyResults = &cobra.Command{
	Use:   "collect-survey-results",
	Short: "Generates JSON format required by the tool by inspecting a google form's responses",
	Args: func(cmd *cobra.Command, args []string) error {
		if surveyFilepath == "" {
			return fmt.Errorf("survey required")
		}

		var err error

		if emailAddressColumn < 0 {
			return fmt.Errorf("invalid email-address")
		}

		if platformsColumn < 0 {
			return fmt.Errorf("invalid platform")
		}

		if ignoreColumnsStr != "" {
			ignoreColumnsStrs := strings.Split(ignoreColumnsStr, ",")
			ignoreColumns = make([]int, len(ignoreColumnsStrs))
			for index, ignoreColumnStr := range ignoreColumnsStrs {
				ignoreColumns[index], err = strconv.Atoi(ignoreColumnStr)
				if err != nil {
					return fmt.Errorf("invalid ignore")
				}
			}
		}

		if previousParticipantsFilepath != "" {
			previousParticipants, err = common.GetParticipantsFromJSONFile(previousParticipantsFilepath, false)
			if err != nil {
				return err
			}
		}

		newParticipants, err = common.GetParticipantsFromCSVFile(surveyFilepath, emailAddressColumn, platformsColumn, ignoreColumns, platformsSeparator)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		participants := common.MergeParticipants(newParticipants, previousParticipants)
		json, err := common.GenerateParticipantsJSON(participants)
		if err != nil {
			return err
		}

		fmt.Println(json)
		return nil
	},
}
