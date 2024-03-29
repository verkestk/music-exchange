package participant

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/verkestk/music-exchange/src/email"
)

// Pair is a pair of participants
type Pair struct {
	Giver    *Participant
	Receiver *Participant
}

// Response is a simple question/answer pair
type Response struct {
	Question string
	Answer   string
}

// Participant represents a person who will give/receive in the exchange
type Participant struct {
	EmailAddress     string // must be unique
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

// GetParticipantsFromJSON transforms the json to a slice of Particpant
func GetParticipantsFromJSON(jsonStr string, skip bool) ([]*Participant, error) {
	participants := []*Participant{}
	allParticipants := []*Participant{}
	err := json.Unmarshal([]byte(jsonStr), &allParticipants)

	if err != nil {
		return nil, fmt.Errorf("error decoding participants json: %w", err)
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

// GetParticipantsFromCSV transforms the CSV to a slice of Particpant
func GetParticipantsFromCSV(csvStr string, emailAddressColumn, platformsColumn int, ignoreColumns []int, platformsSeparator string) ([]*Participant, error) {
	r := csv.NewReader(strings.NewReader(csvStr))
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
		if index == emailAddressColumn || index == platformsColumn {
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
			EmailAddress: row[emailAddressColumn],
			Platforms:    strings.Split(row[platformsColumn], platformsSeparator),
			Responses:    []*Response{},
		}

		// get the questions
		for index, field := range row {
			if index == emailAddressColumn || index == platformsColumn {
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

func (p *Participant) getInstructions(recipient *Participant, instructionsTMPL *template.Template) (string, error) {
	recipientName := recipient.EmailAddress

	instr := &Instructions{
		GiverName:         p.EmailAddress,
		ReceiverName:      recipientName,
		ReceiverResponses: recipient.Responses,
		ReceiverPlatforms: recipient.Platforms,
	}

	tmplBytes := bytes.Buffer{}
	err := instructionsTMPL.Execute(&tmplBytes, instr)
	if err != nil {
		return "", fmt.Errorf("error creating instructions text: %w", err)
	}

	return tmplBytes.String(), nil
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
	if previousParticipants == nil || len(previousParticipants) == 0 {
		return participants
	}

	// map participants by ID
	currentParticipantsMap := map[string]*Participant{}
	for _, participant := range participants {
		currentParticipantsMap[participant.EmailAddress] = participant
	}

	for _, previousParticipant := range previousParticipants {
		if currentParticipant, ok := currentParticipantsMap[previousParticipant.EmailAddress]; ok {
			// this previous participant is already in the particpants list
			currentParticipant.LatestRecipients = previousParticipant.LatestRecipients
		} else {
			// this previous participant isn't participating this time
			participants = append(participants, &Participant{
				EmailAddress:     previousParticipant.EmailAddress,
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

// Score provides a score for the pair based on previous pairings
func (p *Pair) Score() float64 {
	score := float64(0)
	for i, previousRecipient := range p.Giver.LatestRecipients {
		if previousRecipient == p.Receiver.EmailAddress {
			score += float64(1) / float64(i+1)
		}
	}

	return score
}

// IsRepeat returns true if this is the same recipient as the most recent exchange
func (p *Pair) IsRepeat() bool {
	// TODO
	return len(p.Giver.LatestRecipients) > 0 && p.Giver.LatestRecipients[0] == p.Receiver.EmailAddress
}

// EmailInstructions emails the instructions to the participants based on provided email template
func EmailInstructions(pairs []*Pair, tmpl *template.Template, subject, emailTestRecipient string, sender email.Sender) error {
	// TODO how to handling errors?
	for _, pair := range pairs {
		instructions, err := pair.Giver.getInstructions(pair.Receiver, tmpl)
		if err != nil {
			return err
		}

		instructionsRecipient := pair.Giver.EmailAddress
		if emailTestRecipient != "" {
			instructionsRecipient = emailTestRecipient
		}
		fmt.Printf("sending instructions for %s to %s\n", pair.Giver.EmailAddress, instructionsRecipient)
		err = sender.SendMail(subject, instructions, instructionsRecipient)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteInstructions write the instructions for the participants to local files based on provided template
func WriteInstructions(pairs []*Pair, tmpl *template.Template, extension string) error {
	for _, pair := range pairs {
		instructions, err := pair.Giver.getInstructions(pair.Receiver, tmpl)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s.%s", pair.Giver.EmailAddress, extension), []byte(instructions), 0644)
		if err != nil {
			return fmt.Errorf("error writing instructions file: %w", err)
		}
	}
	return nil
}

// UpdateLatestRecipients adds in the latest pairing assignments to the participants
func UpdateLatestRecipients(participants []*Participant, pairs []*Pair) {
	// capture the new assignments in a map
	newAssignmentsGiverToRecipient := map[string]string{}
	for _, pair := range pairs {
		newAssignmentsGiverToRecipient[pair.Giver.EmailAddress] = pair.Receiver.EmailAddress
	}

	// prepend the LatestRecipients with the new assignments
	for _, participant := range participants {
		if recipient, ok := newAssignmentsGiverToRecipient[participant.EmailAddress]; ok {
			participant.LatestRecipients = append([]string{recipient}, participant.LatestRecipients...)
		}
	}
}

// UpdateParticipantsJSON takes
func UpdateParticipantsJSON(participantsFilepath string, participants []*Participant) error {
	// generate the JSON
	jsonStr, err := GenerateParticipantsJSON(participants)
	if err != nil {
		return err
	}

	// overwrite the original JSON file
	err = ioutil.WriteFile(participantsFilepath, []byte(jsonStr), 0644)
	if err != nil {
		return fmt.Errorf("error overwriting json file: %w", err)
	}

	return nil
}
