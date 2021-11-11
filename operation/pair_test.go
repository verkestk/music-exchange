package operation

import (
	"testing"
)

const (
	validPairTestJSON = `[{
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

	validTMPL = `{{.GiverName}} to {{.ReceiverName}}

{{ range .ReceiverResponses }}
{{ .Question }}
{{ .Answer }}
{{ end }}

{{ range .ReceiverPlatforms }}- {{ . }}</li>
{{ end }}`
)

func getValidPairConfig() *PairConfig {
	return &PairConfig{
		ParticipantsJSON:        validPairTestJSON,
		InstructionsTemplateStr: validTMPL,
		Algorithm:               BFScored,
	}
}

func Test_PairConfig_Prepare_ParticipantsOptions(t *testing.T) {
	config := getValidPairConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if len(config.participants) != 3 {
		t.Errorf("expected participants to have length 3 - got %d", len(config.participants))
	}
	if len(config.filteredParticipants) != 2 {
		t.Errorf("expected filteredParticipants to have length 2 - got %d", len(config.filteredParticipants))
	}

	// missing participants
	config = getValidPairConfig()
	config.ParticipantsJSON = ""
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.participants) != 0 {
		t.Errorf("expected participants to be empty")
	}
	if len(config.filteredParticipants) != 0 {
		t.Errorf("expected filteredParticipants to be empty")
	}

	// invalid participants
	config = getValidPairConfig()
	config.ParticipantsJSON = "THIS IS NOT JSON <<<<<< }}} []]]{"
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.participants) != 0 {
		t.Errorf("expected participants to be empty")
	}
	if len(config.filteredParticipants) != 0 {
		t.Errorf("expected filteredParticipants to be empty")
	}
}
func Test_PairConfig_Prepare_InstructionsOptions(t *testing.T) {
	config := getValidPairConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if config.instructionsTMPL == nil {
		t.Errorf("instructions template should not be nil")
	}

	// missing participants
	config = getValidPairConfig()
	config.InstructionsTemplateStr = ""
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if config.instructionsTMPL != nil {
		t.Errorf("instructions template should be nil")
	}

	// invalid participants
	config = getValidPairConfig()
	config.InstructionsTemplateStr = "{ ..... { { } THIS IS NOT VALID <<<<<< }}} []]]{{{{{{"
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}
	if len(config.participants) != 0 {
		t.Errorf("expected participants to be empty")
	}
}
func Test_PairConfig_Prepare_AlgorithmOptions(t *testing.T) {
	config := getValidPairConfig()
	err := config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}

	// if Avoid is already set, AllowRepeatParticipants shouldn't impact it
	config = getValidPairConfig()
	config.Avoid = 2
	config.AllowRepeatParticipants = true
	err = config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if config.Avoid != 2 {
		t.Errorf("expect avoid value of 2 - got %d", config.Avoid)
	}

	// if Avoid is not already set, AllowRepeatParticipants should impact it
	config = getValidPairConfig()
	config.Avoid = 0 // equivalent of unset
	config.AllowRepeatParticipants = true
	err = config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}
	if config.Avoid != 1 {
		t.Errorf("expect avoid value of 1 - got %d", config.Avoid)
	}

	// valid algorith 1
	config = getValidPairConfig()
	config.Algorithm = BFRandom
	err = config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}

	// valid algorithm 2
	config = getValidPairConfig()
	config.Algorithm = BFScored
	err = config.Prepare()
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if !config.prepared {
		t.Errorf("config should be prepared")
	}

	// invalid algorithms
	config = getValidPairConfig()
	config.Algorithm = 100 // invalid
	err = config.Prepare()
	if err == nil {
		t.Errorf("expected error")
	}
	if config.prepared {
		t.Errorf("config should not be prepared")
	}

}
