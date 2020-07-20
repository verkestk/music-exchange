package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"text/template"
	"time"
)

type response struct {
	Question string
	Answer   string
}

type participant struct {
	Name             string
	ID               string
	LatestRecipients []string
	Platforms        []string
	Responses        []*response
}

type instructions struct {
	GiverName         string
	ReceiverName      string
	ReceiverResponses []*response
	ReceiverPlatforms []string
}

var (
	participantFilepath = flag.String("filepath", "", "path to the JSON participants file")
	avoid               = flag.Int("avoid", 0, "how many recipients back to avoid")
	mdTemplate          = `# Secret Music Exchange!

Dear {{.GiverName}},

You will be selecting music for **{{.ReceiverName}}**.

Here's what they had to say about themselves:

<dl>
{{ range .ReceiverResponses }}
  <dt>{{ .Question }}</dt>
  <dd>{{ .Answer }}</dd>
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
)

func main() {
	flag.Parse()

	if participantFilepath == nil || *participantFilepath == "" {
		flag.Usage()
		log.Fatalln("filepath required")
	}

	tmpl := template.Must(template.New("mdTemplate").Parse(mdTemplate))

	byteValue, err := ioutil.ReadFile(*participantFilepath)
	if err != nil {
		log.Fatalln(err)
	}

	participants := []*participant{}
	err = json.Unmarshal(byteValue, &participants)

	if err != nil {
		log.Fatalln(err)
	}

	if len(participants) < 2 {
		log.Fatalf("paricipant count: %d, minimum required: %d\n", len(participants), 2)
	}

	targets := copy(participants)
	shuffle(participants, targets)

	fmt.Println("participants")
	for i := range participants {

		instr := &instructions{
			GiverName:         participants[i].Name,
			ReceiverName:      targets[i].Name,
			ReceiverResponses: targets[i].Responses,
			ReceiverPlatforms: targets[i].Platforms,
		}

		tmplBytes := bytes.Buffer{}
		err := tmpl.Execute(&tmplBytes, instr)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s.md", participants[i].ID), tmplBytes.Bytes(), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func copy(input []*participant) []*participant {
	output := make([]*participant, len(input))
	for i := range input {
		output[i] = input[i]
	}

	return output
}

func shuffle(givers, receivers []*participant) {
	rand.Seed(time.Now().UnixNano())

	for !ok(givers, receivers) {
		rand.Shuffle(len(receivers), func(i, j int) { receivers[i], receivers[j] = receivers[j], receivers[i] })
	}
}

func ok(givers, receivers []*participant) bool {
	for i := range givers {

		// the giver and the receiver can't be the same person
		if givers[i].ID == receivers[i].ID {
			log.Printf("%s is the same as %s\n", givers[i].ID, receivers[i].ID)
			return false
		}

		// avoid pairing the same people
		for j := 0; j < *avoid && j < len(givers[i].LatestRecipients); j++ {
			if givers[i].LatestRecipients[j] == receivers[i].ID {
				log.Printf("%s gave to %s %d times ago\n", givers[i].ID, receivers[i].ID, j+1)
				return false
			}
		}

		// the giver and the receiver must have at least one platform in common
		if countOverlap(givers[i].Platforms, receivers[i].Platforms) == 0 {
			log.Printf("%s and %s have no platforms in common\n", givers[i].ID, receivers[i].ID)
			return false
		}
	}

	return true
}

func countOverlap(s1, s2 []string) int {
	overlap := map[string]bool{}
	for _, s1Val := range s1 {
		for _, s2Val := range s2 {
			if s1Val == s2Val {
				overlap[s1Val] = true
			}
		}
	}

	return len(overlap)
}
