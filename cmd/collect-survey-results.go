package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/verkestk/music-exchange/operation"
)

var (
	collectConfig = &operation.CollectConfig{}
)

func init() {
	rootCmd.AddCommand(collectSuveyResults)
	collectSuveyResults.Flags().StringVarP(&collectConfig.SurveyFilepath, "survey", "s", "", "filepath of CSV survey results")
	collectSuveyResults.Flags().StringVarP(&collectConfig.PreviousParticipantsFilepath, "previous-participants", "p", "", "filepath of the JSON file used for the last exchange")
	collectSuveyResults.Flags().IntVarP(&collectConfig.EmailAddressColumn, "email-address", "e", -1, "index of the column containing the email address")
	collectSuveyResults.Flags().IntVarP(&collectConfig.PlatformsColumn, "platforms", "l", -1, "index of the column containing the music platforms")
	collectSuveyResults.Flags().StringVarP(&collectConfig.IgnoreColumnsStr, "ignore", "i", "", "indexes of columns to ignore")
	collectSuveyResults.Flags().StringVarP(&collectConfig.PlatformsSeparator, "separator", "a", ";", "separator character between platform choices")
	collectSuveyResults.MarkFlagRequired("survey")
	collectSuveyResults.MarkFlagRequired("email-address")
	collectSuveyResults.MarkFlagRequired("platformsColumn")
}

var collectSuveyResults = &cobra.Command{
	Use:   "collect-survey-results",
	Short: "Generates JSON format required by the tool by inspecting a google form's responses",
	Args: func(cmd *cobra.Command, args []string) error {
		return collectConfig.Prepare()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		json, err := operation.DoCollect(collectConfig)
		if err != nil {
			return err
		}

		fmt.Println(json)
		return nil
	},
}
