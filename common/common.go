package common

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// Response is a simple question/answer pair
type Response struct {
	Question string
	Answer   string
}

// Participant represents a person who will give/receive in the exchange
type Participant struct {
	Username         string // must be unique
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

// GetParticipantsFromJSONFile transforms the json in a file to a slice of Particpant
func GetParticipantsFromJSONFile(filepath string, skip bool) ([]*Participant, error) {
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
		if !p.Skip || !skip {
			participants = append(participants, p)
		}
	}

	if len(participants) < 2 {
		return nil, fmt.Errorf("paricipant count: %d, minimum required: %d", len(participants), 2)
	}

	return participants, nil
}

// GetParticipantsFromCSVFile transforms the CSV in a file to a slice of Particpant
func GetParticipantsFromCSVFile(filepath string, usernameColumn, platformsColumn int, ignoreColumns []int, platformsSeparator string) ([]*Participant, error) {
	reader, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(reader)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 3 {
		return nil, fmt.Errorf("csv file must have at least 3 rows - 1 header and 2+ participants")
	}

	// get the list of arbitrary questions from the header row
	questions := []string{}
	for index, questionText := range rows[0] {
		if index == usernameColumn || index == platformsColumn {
			continue
		}
		ignore := false
		for _, ignoreColumn := range ignoreColumns {
			if index == ignoreColumn {
				ignore = true
				break
			}
		}
		if ignore {
			continue
		}

		questions = append(questions, questionText)
	}

	// get the participants from the rest of the columns
	participants := []*Participant{}
	for _, row := range rows[1:] {
		participant := &Participant{
			Username:  row[usernameColumn],
			Platforms: strings.Split(row[platformsColumn], platformsSeparator),
			Responses: []*Response{},
		}

		// get the questions
		for index, field := range row {
			if index == usernameColumn || index == platformsColumn {
				continue
			}
			ignore := false
			for _, ignoreColumn := range ignoreColumns {
				if index == ignoreColumn {
					ignore = true
					break
				}
			}
			if ignore {
				continue
			}
			participant.Responses = append(participant.Responses, &Response{
				Question: questions[len(participant.Responses)],
				Answer:   field,
			})
		}

		participants = append(participants, participant)
	}

	return participants, nil
}

// WriteInstructions outputs the instructions as MD in a file
func (p *Participant) WriteInstructions(recipient *Participant, instructionsTMPL *template.Template) error {
	recipientName := recipient.Username

	instr := &Instructions{
		GiverName:         p.Username,
		ReceiverName:      recipientName,
		ReceiverResponses: recipient.Responses,
		ReceiverPlatforms: recipient.Platforms,
	}

	tmplBytes := bytes.Buffer{}
	err := instructionsTMPL.Execute(&tmplBytes, instr)
	if err != nil {
		return fmt.Errorf("error creating instructions text: %w", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.md", p.Username), tmplBytes.Bytes(), 0644)
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

// MergeParticipants takes a set of participants and merges in previous participant data
// Previous participants who also exist in the participants slice have their pairing history included.
// Previous participants who don't exist in the participants slice are copied in with the pairing history and marked to be skipped.
func MergeParticipants(participants, previousParticipants []*Participant) []*Participant {
	// map participants by ID
	currentParticipantsMap := map[string]*Participant{}
	for _, participant := range participants {
		currentParticipantsMap[participant.Username] = participant
	}

	for _, previousParticipant := range previousParticipants {
		if currentParticipant, ok := currentParticipantsMap[previousParticipant.Username]; ok {
			// this previous participant is already in the particpants list
			currentParticipant.LatestRecipients = previousParticipant.LatestRecipients
		} else {
			// this previous participant isn't participating this time
			participants = append(participants, &Participant{
				Username:         previousParticipant.Username,
				LatestRecipients: previousParticipant.LatestRecipients,
				Skip:             true,
			})
		}
	}
	return participants
}

// GenerateParticipantsJSON generates readable JSON for the participants
func GenerateParticipantsJSON(participants []*Participant) (string, error) {
	bytes, err := json.MarshalIndent(participants, "", "  ")
	return string(bytes), err
}
