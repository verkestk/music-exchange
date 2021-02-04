package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
)

// Response is a simple question/answer pair
type Response struct {
	Question string
	Answer   string
}

// Participant represents a person who will give/receive in the exchange
type Participant struct {
	Name             string
	ID               string
	Skip             bool
	LatestRecipients []string
	Platforms        []string
	Responses        []*Response
}

// Instructions contains the information necessary for generating exchange instructions
// Works with the instructions MD template
type Instructions struct {
	GiverName         string
	ReceiverName      string
	ReceiverResponses []*Response
	ReceiverPlatforms []string
}

// GetParticipantsFromFile transforms the json is a file to a slice of Particpant
func GetParticipantsFromFile(filepath string) ([]*Participant, error) {
	byteValue, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading from file path %s: %w", filepath, err)
	}

	participants := []*Participant{}
	allParticipants := []*Participant{}
	err = json.Unmarshal(byteValue, &allParticipants)

	if err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	for _, p := range allParticipants {
		if !p.Skip {
			participants = append(participants, p)
		}
	}

	if len(participants) < 2 {
		return nil, fmt.Errorf("paricipant count: %d, minimum required: %d", len(participants), 2)
	}

	return participants, nil
}

// WriteInstructions outputs the instructions as MD in a file
func (p *Participant) WriteInstructions(recipient *Participant, instructionsTMPL *template.Template) error {
	instr := &Instructions{
		GiverName:         p.Name,
		ReceiverName:      recipient.Name,
		ReceiverResponses: recipient.Responses,
		ReceiverPlatforms: recipient.Platforms,
	}

	tmplBytes := bytes.Buffer{}
	err := instructionsTMPL.Execute(&tmplBytes, instr)
	if err != nil {
		return fmt.Errorf("error creating instructions text: %w", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.md", p.ID), tmplBytes.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing instructions file: %w", err)
	}

	return nil
}

// IsCompatible returns true if the participants have at least one music platform in common
func (p *Participant) IsCompatible(otherParticipant *Participant) bool {
	for _, pPlatform := range p.Platforms {
		for _, opPlatform := range otherParticipant.Platforms {
			if pPlatform == opPlatform {
				return true
			}
		}
	}

	return false
}
