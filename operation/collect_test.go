package operation

import (
	"testing"
)

const (
	emailAddressColumn = 0
	platformsColumn    = 3
	ignoreColumnsStr   = "4,5"
	platformsSeparator = ";"

	validCSV = `Email Address,What music do you like?,Who's your favorite musician?,What music platforms do you use?,Skip,Skip
woody@toy.story,I like a good western toon.,Randy Newman,Spotify;YouTube,Ignore,Ignore
buzz@toy.story,Space Rock,Randy Newman,Spotify;YouTube,Ignore,Ignore
rex@toy.story,Prehistoric,Randy Newman,Apple Music;YouTube;YouTube Music;Pandora,Ignore,Ignore
hamm@toy.story,Classy Toons,Randy Newman,Spotify;Pandora,Ignore,Ignore
jessie@toy.story,good ol' classics from the wild west.,Randy Newman,Spotify;Pandora,Ignore,Ignore`

	validCollectTestJSON = `[{
  "EmailAddress": "woody@toy.story",
  "Skip": true,
  "LatestRecipients": ["buzz@toy.story"],
  "Platforms": [
    "Spotify"
  ],
  "Responses": [{
    "Question": "What music do you like?",
    "Answer": "songs that make me feel like a real cowboy"
  },{
    "Question": "Who's your favorite musician?",
    "Answer": "Randy Newman"
  }]
}, {
  "EmailAddress": "buzz@toy.story",
  "LatestRecipients": ["slinky@toy.story", "woody@toy.story"],
  "Platforms": [
    "Spotify"
  ],
  "Responses": [{
    "Question": "What music do you like?",
    "Answer": "infinitely good music"
  },{
    "Question": "Who's your favorite musician?",
    "Answer": "Randy Newman"
  }]
}, {
  "EmailAddress": "slinky@toy.story",
  "LatestRecipients": ["buzz@toy.story"],
  "Platforms": [
    ""
  ],
  "Responses": [{
    "Question": "What music do you like?",
    "Answer": "long songs"
  },{
    "Question": "Who's your favorite musician?",
    "Answer": "Randy Newman"
  }]
}]`
)

func getValidCollectConfig() *CollectConfig {
	return &CollectConfig{
		SurveyCSV:                validCSV,
		PreviousParticipantsJSON: validCollectTestJSON,
		EmailAddressColumn:       emailAddressColumn,
		PlatformsColumn:          platformsColumn,
		IgnoreColumnsStr:         ignoreColumnsStr,
		PlatformsSeparator:       platformsSeparator,
	}
}

func Test_CollectConfig_Prepare_SurveyContentOptions(t *testing.T) {
	config := getValidCollectConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if len(config.newParticipants) == 0 {
		t.Errorf("newParticipants should not be empty")
	}

	// missing survey
	config = getValidCollectConfig()
	config.SurveyCSV = ""
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.newParticipants) != 0 {
		t.Errorf("newParticipants should not empty")
	}

	// invalid CSV
	config = getValidCollectConfig()
	config.SurveyCSV = validCollectTestJSON // AKA invalid CSV
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.newParticipants) != 0 {
		t.Errorf("newParticipants should not empty")
	}
}

func Test_CollectConfig_Prepare_PreviousParticipantsContentOptions(t *testing.T) {
	config := getValidCollectConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if len(config.previousParticipants) == 0 {
		t.Errorf("previousParticipants should not be empty")
	}

	// missing JSON (still valid)
	config = getValidCollectConfig()
	config.PreviousParticipantsJSON = ""
	err = config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if len(config.previousParticipants) != 0 {
		t.Errorf("previousParticipants should be empty")
	}

	// invalid CSV
	config = getValidCollectConfig()
	config.PreviousParticipantsJSON = validCSV // AKA invalid JSON
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.previousParticipants) != 0 {
		t.Errorf("previousParticipants should not empty")
	}
}

func Test_CollectConfig_Prepare_CSVParsingOptions(t *testing.T) {
	config := getValidCollectConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if len(config.previousParticipants) == 0 {
		t.Errorf("previousParticipants should not be empty")
	}
	if len(config.ignoreColumns) != 2 {
		t.Errorf("expected 2 ignored columns - got %d", len(config.ignoreColumns))
	}
	if config.ignoreColumns[0] != 4 {
		t.Errorf("expected ignored column of 4 - got %d", config.ignoreColumns[0])
	}
	if config.ignoreColumns[1] != 5 {
		t.Errorf("expected ignored column of 5 - got %d", config.ignoreColumns[1])
	}

	// invalid email address column
	config = getValidCollectConfig()
	config.EmailAddressColumn = -1
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.newParticipants) != 0 {
		t.Errorf("newParticipants should be empty")
	}

	// invalid platforms column
	config = getValidCollectConfig()
	config.PlatformsColumn = -1
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.newParticipants) != 0 {
		t.Errorf("newParticipants should be empty")
	}

	// missing ignore columns (still valid)
	config = getValidCollectConfig()
	config.IgnoreColumnsStr = ""
	err = config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if len(config.previousParticipants) == 0 {
		t.Errorf("previousParticipants should not be empty")
	}

	// invalid ignore columns
	config = getValidCollectConfig()
	config.IgnoreColumnsStr = "ABC"
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.newParticipants) != 0 {
		t.Errorf("newParticipants should be empty")
	}

	// missing platforms separator
	config = getValidCollectConfig()
	config.PlatformsSeparator = ""
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.newParticipants) != 0 {
		t.Errorf("newParticipants should be empty")
	}
}

func Test_DoCollect(t *testing.T) {
	config := getValidCollectConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	newParticipantsJSON, err := DoCollect(config)
	if newParticipantsJSON == "" {
		t.Errorf("new participants JSON should not be empty")
	}
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}

	// skipping prepare
	config = getValidCollectConfig()
	newParticipantsJSON, err = DoCollect(config)
	if newParticipantsJSON != "" {
		t.Errorf("new participants JSON should be empty")
	}
	if err == nil {
		t.Errorf("expected error")
	}
}
