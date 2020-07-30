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
// Works with InstructionsMdTemplate
type Instructions struct {
	GiverName         string
	ReceiverName      string
	ReceiverResponses []*Response
	ReceiverPlatforms []string
}

// InstructionsMdTemplate is a template for exchange instructions
// works with Instructions struct
var instructionsMdTemplateStr = `# Secret Music Exchange!

Dear {{.GiverName}},

You will be selecting music for **{{.ReceiverName}}**.

Here's what they had to say about themselves:

<dl>
{{ range .ReceiverResponses }}
  <dt>{{ .Question }}</dt>
  <dd>{{ .Answer }}</dd><br/>
{{ end }}
</dl>

### Instructions

##### 1: Select Music
Create a sharable playlist (between 45 and 75 minutes) on one of the platforms your recipient uses:

{{ range .ReceiverPlatforms }}- {{ . }}
{{ end }}
##### 2: Copy Link
If you are using an album, copy the album link from Spotify. If you are creating a playlist, copy the playlist link.

##### 3: Share Music
We'll do this all at once. When it's time, you will share that link on Slack. @mention then in the #music-appreciation channel, so we can all see what you picked.

##### 4: Enjoy the music shared with you!
You'll be getting some music to enjoy too.
`

var instructionsMdTemplate *template.Template

func init() {
	instructionsMdTemplate = template.Must(template.New("mdTemplate").Parse(instructionsMdTemplateStr))
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
func (p *Participant) WriteInstructions(recipient *Participant) error {
	instr := &Instructions{
		GiverName:         p.Name,
		ReceiverName:      recipient.Name,
		ReceiverResponses: recipient.Responses,
		ReceiverPlatforms: recipient.Platforms,
	}

	tmplBytes := bytes.Buffer{}
	err := instructionsMdTemplate.Execute(&tmplBytes, instr)
	if err != nil {
		return fmt.Errorf("error creating instructions text: %w", err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.md", p.ID), tmplBytes.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing instructions file: %w", err)
	}

	return nil
}
