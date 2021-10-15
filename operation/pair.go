package operation

import (
	"fmt"
	"text/template"

	"github.com/verkestk/music-exchange/src/algorithms/bfrandom"
	"github.com/verkestk/music-exchange/src/algorithms/bfscored"
	"github.com/verkestk/music-exchange/src/email"
	"github.com/verkestk/music-exchange/src/participant"
)

// indentifiers for pairing algorithms
const (
	BFRandom = iota
	BFScored
)

// PairConfig contains the inputs neccessary to pair the participants
type PairConfig struct {
	// General Configuration
	ParticipantsFilepath   string
	InstructionsFilepath   string
	UpdateParticipantsFile bool
	Algorithm              int // see constants

	// For Writing Instructions to File
	InstructionsFileExtension string

	// For Emailing Instructions
	EmailInstructions  bool
	EmailSubject       string
	EmailTestRecipient string
	SMTPHostEnvVar     string
	SMTPPortEnvVar     string
	SMTPUsernameEnvVar string
	SMTPPasswordEnvVar string

	// Only for the BF-Random algorithm
	Avoid int

	participants     []*participant.Participant
	instructionsTMPL *template.Template
}

// Prepare intakes the configuration, processes and validates
func (config *PairConfig) Prepare() error {
	if config.ParticipantsFilepath == "" {
		return fmt.Errorf("participants required")
	}

	if config.InstructionsFilepath == "" {
		return fmt.Errorf("instructions required")
	}

	var err error
	config.participants, err = participant.GetParticipantsFromJSONFile(config.ParticipantsFilepath, true)
	if err != nil {
		return err
	}

	config.instructionsTMPL, err = template.ParseFiles(config.InstructionsFilepath)
	if err != nil {
		return err
	}

	return nil
}

// DoPair performs the operation
func DoPair(config *PairConfig) error {

	var pairs []*participant.Pair
	var err error

	switch config.Algorithm {
	case BFRandom:
		pairs, err = bfrandom.DoExchange(config.participants, config.Avoid)
	case BFScored:
		pairs, err = bfscored.DoExchange(config.participants)
	}
	if err != nil {
		return err
	}

	if config.EmailInstructions {
		sender := email.GetGmailSender(config.SMTPHostEnvVar, config.SMTPPortEnvVar, config.SMTPUsernameEnvVar, config.SMTPPasswordEnvVar)
		err = participant.EmailInstructions(pairs, config.instructionsTMPL, config.EmailSubject, config.EmailTestRecipient, sender)
	} else {
		err = participant.WriteInstructions(pairs, config.instructionsTMPL, config.InstructionsFileExtension)
	}

	if err != nil {
		return err
	}

	return participant.UpdateParticipantsJSON(config.ParticipantsFilepath, pairs)
}
