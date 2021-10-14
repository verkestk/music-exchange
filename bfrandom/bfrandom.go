package bfrandom

import (
	"log"
	"math/rand"
	"text/template"
	"time"

	"github.com/verkestk/music-exchange/common"
)

// DoExchange matches particpants as givers and recipients, generating files with instructions for each participant
func DoExchange(participants []*common.Participant, instructionsTMPL *template.Template, avoid int) ([]*common.Pair, error) {

	targets := copy(participants)
	shuffle(participants, targets, avoid)

	pairs := make([]*common.Pair, len(participants))
	for index := range participants {
		pairs[index] = &common.Pair{Giver: participants[index], Receiver: targets[index]}
	}

	return pairs, nil
}

func copy(input []*common.Participant) []*common.Participant {
	output := make([]*common.Participant, len(input))
	for i := range input {
		output[i] = input[i]
	}

	return output
}

func shuffle(givers, receivers []*common.Participant, avoid int) {
	rand.Seed(time.Now().UnixNano())

	for !ok(givers, receivers, avoid) {
		rand.Shuffle(len(receivers), func(i, j int) { receivers[i], receivers[j] = receivers[j], receivers[i] })
	}
}

func ok(givers, receivers []*common.Participant, avoid int) bool {
	for i := range givers {

		// the giver and the receiver can't be the same person
		if givers[i].EmailAddress == receivers[i].EmailAddress {
			log.Printf("%s is the same as %s\n", givers[i].EmailAddress, receivers[i].EmailAddress)
			return false
		}

		// avoid pairing the same people
		for j := 0; j < avoid && j < len(givers[i].LatestRecipients); j++ {
			if givers[i].LatestRecipients[j] == receivers[i].EmailAddress {
				log.Printf("%s gave to %s %d times ago\n", givers[i].EmailAddress, receivers[i].EmailAddress, j+1)
				return false
			}
		}

		// the giver and the receiver must have at least one platform in common
		if !givers[i].IsCompatible(receivers[i]) {
			log.Printf("%s and %s have no platforms in common\n", givers[i].EmailAddress, receivers[i].EmailAddress)
			return false
		}
	}

	return true
}
