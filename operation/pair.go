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
	AllowRepeatParticipants bool
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

	participants         []*participant.Participant // includes skipped participants
	filteredParticipants []*participant.Participant // does not include skipped participants
	instructionsTMPL     *template.Template

	prepared bool
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

	if config.AllowRepeatParticipants && config.Avoid == 0 {
		config.Avoid = 1
	}

	if config.UpdateParticipantsFile && config.ParticipantsFilepath == "" {
		return fmt.Errorf("participants filepath required to update the file")
	}

	if config.ParticipantsFilepath == "" && config.ParticipantsJSON == "" {
		return fmt.Errorf("participants required")
	}
	if config.ParticipantsJSON == "" {
		// generate JSON from file
		byteValue, err := ioutil.ReadFile(config.ParticipantsFilepath)
		if err != nil {
			return fmt.Errorf("cannot perform pairing: error reading from participants file path %s: %w", config.ParticipantsFilepath, err)
		}
		config.ParticipantsJSON = string(byteValue)
	}

	config.participants, err = participant.GetParticipantsFromJSON(config.ParticipantsJSON, false)
	if err != nil {
		return fmt.Errorf("cannot perform pairing: error getting participants: %w", err)
	}

	config.filteredParticipants = []*participant.Participant{}
	for _, p := range config.participants {
		if !p.Skip {
			config.filteredParticipants = append(config.filteredParticipants, p)
		}
	}

	if config.Algorithm != BFRandom && config.Algorithm != BFScored {
		return fmt.Errorf("unsupported algorithm")
	}

	config.prepared = true
	return nil
}

// DoPair performs the operation
func DoPair(config *PairConfig) (string, error) {
	if !config.prepared {
		return "", fmt.Errorf("pair config not prepared")
	}

	var pairs []*participant.Pair
	var err error

	switch config.Algorithm {
	case BFRandom:
		pairs, err = bfrandom.DoExchange(config.filteredParticipants, config.Avoid)
	case BFScored:
		pairs, err = bfscored.DoExchange(config.filteredParticipants, config.AllowRepeatParticipants)
	default:
		return "", fmt.Errorf("unsupported algorithm: %d", config.Algorithm)
	}
	if err != nil {
		return "", err
	}

	if config.EmailInstructions {
		err = participant.EmailInstructions(pairs, config.instructionsTMPL, config.EmailSubject, config.EmailTestRecipient, config.EmailSender)
	} else {
		err = participant.WriteInstructions(pairs, config.instructionsTMPL, config.InstructionsFileExtension)
	}

	if err != nil {
		return "", err
	}

	participant.UpdateLatestRecipients(config.participants, pairs)

	if config.UpdateParticipantsFile {
		err = participant.UpdateParticipantsJSON(config.ParticipantsFilepath, config.participants)
		if err != nil {
			return "", fmt.Errorf("error writing participants JSON: %w", err)
		}
	}

	return participant.GenerateParticipantsJSON(config.participants)
}
