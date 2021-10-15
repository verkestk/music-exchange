package operation

import (
	"fmt"
	"io/ioutil"
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
	ParticipantsFilepath    string
	ParticipantsJSON        string
	InstructionsFilepath    string
	InstructionsTemplateStr string
	UpdateParticipantsFile  bool
	Algorithm               int // see constants

	// For Writing Instructions to File
	InstructionsFileExtension string

	// For Emailing Instructions
	EmailInstructions  bool
	EmailSubject       string
	EmailTestRecipient string
	EmailSender        email.Sender

	// Only for the BF-Random algorithm
	Avoid int

	participants     []*participant.Participant
	instructionsTMPL *template.Template
}

// Prepare intakes the configuration, processes and validates
func (config *PairConfig) Prepare() error {
	var err error

	if config.InstructionsFilepath == "" && config.InstructionsTemplateStr == "" {
		return fmt.Errorf("instructions required")
	}
	if config.InstructionsTemplateStr == "" {
		config.instructionsTMPL, err = template.ParseFiles(config.InstructionsFilepath)
		if err != nil {
			return err
		}
	} else {
		config.instructionsTMPL, err = template.New("instructions").Parse(config.InstructionsTemplateStr)
		if err != nil {
			return err
		}
	}

	if config.ParticipantsFilepath == "" && config.ParticipantsJSON == "" {
		return fmt.Errorf("participants required")
	}
	if config.ParticipantsJSON == "" {
		// generate JSON from file
		byteValue, err := ioutil.ReadFile(config.ParticipantsFilepath)
		if err != nil {
			return fmt.Errorf("error reading from file path %s: %w", config.ParticipantsFilepath, err)
		}
		config.ParticipantsJSON = string(byteValue)
	}

	config.participants, err = participant.GetParticipantsFromJSON(config.ParticipantsJSON, true)
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
	default:
		return fmt.Errorf("unsupported algorithm: %d", config.Algorithm)
	}
	if err != nil {
		return err
	}

	if config.EmailInstructions {
		err = participant.EmailInstructions(pairs, config.instructionsTMPL, config.EmailSubject, config.EmailTestRecipient, config.EmailSender)
	} else {
		err = participant.WriteInstructions(pairs, config.instructionsTMPL, config.InstructionsFileExtension)
	}

	if err != nil {
		return err
	}

	return participant.UpdateParticipantsJSON(config.ParticipantsFilepath, pairs)
}
