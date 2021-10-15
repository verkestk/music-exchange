package operation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/verkestk/music-exchange/src/participant"
)

// CollectConfig contains the inputs neccessary for collecting survey results
type CollectConfig struct {
	SurveyFilepath               string
	PreviousParticipantsFilepath string
	EmailAddressColumn           int
	PlatformsColumn              int
	IgnoreColumnsStr             string
	PlatformsSeparator           string

	ignoreColumns        []int
	previousParticipants []*participant.Participant
	newParticipants      []*participant.Participant
}

// Prepare intakes the configuration, processes and validates
func (config *CollectConfig) Prepare() error {

	if config.SurveyFilepath == "" {
		return fmt.Errorf("survey required")
	}

	var err error

	if config.EmailAddressColumn < 0 {
		return fmt.Errorf("invalid email-address")
	}

	if config.PlatformsColumn < 0 {
		return fmt.Errorf("invalid platform")
	}

	if config.IgnoreColumnsStr != "" {
		ignoreColumnsStrs := strings.Split(config.IgnoreColumnsStr, ",")
		config.ignoreColumns = make([]int, len(ignoreColumnsStrs))
		for index, ignoreColumnStr := range ignoreColumnsStrs {
			config.ignoreColumns[index], err = strconv.Atoi(ignoreColumnStr)
			if err != nil {
				return fmt.Errorf("invalid ignore")
			}
		}
	}

	if config.PreviousParticipantsFilepath != "" {
		config.previousParticipants, err = participant.GetParticipantsFromJSONFile(config.PreviousParticipantsFilepath, false)
		if err != nil {
			return err
		}
	}

	config.newParticipants, err = participant.GetParticipantsFromCSVFile(
		config.SurveyFilepath,
		config.EmailAddressColumn,
		config.PlatformsColumn,
		config.ignoreColumns,
		config.PlatformsSeparator,
	)
	if err != nil {
		return err
	}

	return nil
}

// DoCollect performs the operation
func DoCollect(config *CollectConfig) (string, error) {
	participants := participant.MergeParticipants(config.newParticipants, config.previousParticipants)
	return participant.GenerateParticipantsJSON(participants)
}
