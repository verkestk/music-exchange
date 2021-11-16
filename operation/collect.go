package operation

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/verkestk/music-exchange/src/participant"
)

// CollectConfig contains the inputs neccessary for collecting survey results
type CollectConfig struct {
	SurveyFilepath               string
	SurveyCSV                    string
	PreviousParticipantsFilepath string
	PreviousParticipantsJSON     string
	EmailAddressColumn           int
	PlatformsColumn              int
	IgnoreColumnsStr             string // comma-separated
	PlatformsSeparator           string

	ignoreColumns        []int
	previousParticipants []*participant.Participant
	newParticipants      []*participant.Participant

	prepared bool
}

// Prepare intakes the configuration, processes and validates
func (config *CollectConfig) Prepare() error {

	if config.SurveyFilepath == "" && config.SurveyCSV == "" {
		return fmt.Errorf("survey required")
	}

	var err error

	if config.EmailAddressColumn < 0 {
		return fmt.Errorf("invalid email-address")
	}

	if config.PlatformsColumn < 0 {
		return fmt.Errorf("invalid platform")
	}

	if config.PlatformsSeparator == "" {
		return fmt.Errorf("platforms separator required")
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

	if config.PreviousParticipantsJSON == "" && config.PreviousParticipantsFilepath != "" {
		// generate JSON from file
		byteValue, err := ioutil.ReadFile(config.PreviousParticipantsFilepath)
		if err != nil {
			return fmt.Errorf("cannot colect survey results: error reading from participants file path %s: %w", config.PreviousParticipantsFilepath, err)
		}
		config.PreviousParticipantsJSON = string(byteValue)
	}

	if config.PreviousParticipantsJSON != "" {
		config.previousParticipants, err = participant.GetParticipantsFromJSON(config.PreviousParticipantsJSON, false)
		if err != nil {
			return fmt.Errorf("cannot collect survey results: error getting participants: %w", err)
		}
	}

	if config.SurveyCSV == "" {
		// generate JSON from file
		byteValue, err := ioutil.ReadFile(config.SurveyFilepath)
		if err != nil {
			return fmt.Errorf("cannot collect survey results: error reading from survey file path %s: %w", config.SurveyFilepath, err)
		}
		config.SurveyCSV = string(byteValue)
	}
	config.newParticipants, err = participant.GetParticipantsFromCSV(
		config.SurveyCSV,
		config.EmailAddressColumn,
		config.PlatformsColumn,
		config.ignoreColumns,
		config.PlatformsSeparator,
	)
	if err != nil {
		return err
	}

	config.prepared = true
	return nil
}

// DoCollect performs the operation
func DoCollect(config *CollectConfig) (string, error) {
	if !config.prepared {
		return "", fmt.Errorf("collect config not prepared")
	}

	participants := participant.MergeParticipants(config.newParticipants, config.previousParticipants)
	return participant.GenerateParticipantsJSON(participants)
}
