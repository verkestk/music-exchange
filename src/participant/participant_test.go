package participant

import (
	"testing"
)

func Test_IsCompatible(t *testing.T) {
	var (
		woody  = &Participant{Platforms: []string{}}
		buzz   = &Participant{Platforms: []string{"Spotify"}}
		rex    = &Participant{Platforms: []string{"Pandora"}}
		hamm   = &Participant{Platforms: []string{"Pandora", "YouTube"}}
		jessie = &Participant{Platforms: []string{"Spotify", "AmazonMusic"}}
	)

	if woody.IsCompatible(buzz) {
		t.Errorf("woody is compatible with buzz - expected incompatible")
	}

	if woody.IsCompatible(rex) {
		t.Errorf("woody is compatible with rex - expected incompatible")
	}

	if woody.IsCompatible(hamm) {
		t.Errorf("woody is compatible with hamm - expected incompatible")
	}

	if woody.IsCompatible(jessie) {
		t.Errorf("woody is compatible with jessie - expected incompatible")
	}

	if buzz.IsCompatible(woody) {
		t.Errorf("buzz is compatible with woody - expected incompatible")
	}

	if buzz.IsCompatible(rex) {
		t.Errorf("buzz is compatible with rex - expected incompatible")
	}

	if buzz.IsCompatible(hamm) {
		t.Errorf("buzz is compatible with hamm - expected incompatible")
	}

	if !buzz.IsCompatible(jessie) {
		t.Errorf("buzz is incompatible with jessie - expected compatible")
	}

	if rex.IsCompatible(woody) {
		t.Errorf("rex is compatible with woody - expected incompatible")
	}

	if rex.IsCompatible(buzz) {
		t.Errorf("rex is compatible with buzz - expected incompatible")
	}

	if !rex.IsCompatible(hamm) {
		t.Errorf("rex is incompatible with hamm - expected compatible")
	}

	if rex.IsCompatible(jessie) {
		t.Errorf("rex is compatible with jessie - expected incompatible")
	}

	if hamm.IsCompatible(woody) {
		t.Errorf("hamm is compatible with woody - expected incompatible")
	}

	if hamm.IsCompatible(buzz) {
		t.Errorf("hamm is compatible with buzz - expected incompatible")
	}

	if !hamm.IsCompatible(rex) {
		t.Errorf("hamm is incompatible with rex - expected compatible")
	}

	if hamm.IsCompatible(jessie) {
		t.Errorf("hamm is compatible with jessie - expected incompatible")
	}

	if jessie.IsCompatible(woody) {
		t.Errorf("jessie is compatible with woody - expected incompatible")
	}

	if !jessie.IsCompatible(buzz) {
		t.Errorf("jessie is incompatible with buzz - expected compatible")
	}

	if jessie.IsCompatible(rex) {
		t.Errorf("jessie is compatible with rex - expected incompatible")
	}

	if jessie.IsCompatible(hamm) {
		t.Errorf("jessie is compatible with hamm - expected incompatible")
	}
}

func Test_MergeParticipants(t *testing.T) {
	var (
		woody  = &Participant{EmailAddress: "woody@toy.story", Skip: true}
		buzz   = &Participant{EmailAddress: "buzz@toy.story"}
		rex    = &Participant{EmailAddress: "rex@toy.story"}
		hamm   = &Participant{EmailAddress: "hamm@toy.story", Skip: true}
		jessie = &Participant{EmailAddress: "jessie@toy.story"}

		previousParticipants = []*Participant{woody, buzz, rex}
		newParticipants      = []*Participant{rex, hamm, jessie}
	)

	merged := MergeParticipants(newParticipants, nil)
	if len(merged) != 3 {
		t.Errorf("expected 3 participants - got %d", len(merged))
	}
	if merged[0].EmailAddress != rex.EmailAddress {
		t.Errorf("expected rex@toy.story - got %s", merged[0].EmailAddress)
	}
	if merged[1].EmailAddress != hamm.EmailAddress {
		t.Errorf("expected hamm@toy.story - got %s", merged[1].EmailAddress)
	}
	if merged[2].EmailAddress != jessie.EmailAddress {
		t.Errorf("expected jessie@toy.story - got %s", merged[2].EmailAddress)
	}
	if merged[0].Skip {
		t.Errorf("rex should not be skipped")
	}
	if !merged[1].Skip {
		t.Errorf("hamm should be skipped")
	}
	if merged[2].Skip {
		t.Errorf("jessie should not be skipped")
	}

	merged = MergeParticipants(newParticipants, previousParticipants)
	if len(merged) != 5 {
		t.Errorf("expected 5 participants - got %d", len(merged))
	}
	if merged[0].EmailAddress != rex.EmailAddress {
		t.Errorf("expected rex@toy.story - got %s", merged[0].EmailAddress)
	}
	if merged[1].EmailAddress != hamm.EmailAddress {
		t.Errorf("expected hamm@toy.story - got %s", merged[1].EmailAddress)
	}
	if merged[2].EmailAddress != jessie.EmailAddress {
		t.Errorf("expected jessie@toy.story - got %s", merged[2].EmailAddress)
	}
	if merged[3].EmailAddress != woody.EmailAddress {
		t.Errorf("expected woody@toy.story - got %s", merged[2].EmailAddress)
	}
	if merged[4].EmailAddress != buzz.EmailAddress {
		t.Errorf("expected buzz@toy.story - got %s", merged[2].EmailAddress)
	}
	if merged[0].Skip {
		t.Errorf("rex should not be skipped")
	}
	if !merged[1].Skip {
		t.Errorf("hamm should be skipped")
	}
	if merged[2].Skip {
		t.Errorf("jessie should not be skipped")
	}
	if !merged[4].Skip {
		t.Errorf("woody should be skipped")
	}
	if !merged[4].Skip {
		t.Errorf("buzz should not be skipped")
	}
}

func Test_UpdateLatestRecipients(t *testing.T) {
	var (
		woody  = &Participant{EmailAddress: "woody@toy.story", LatestRecipients: []string{"hamm@toy.story"}}
		buzz   = &Participant{EmailAddress: "buzz@toy.story", LatestRecipients: []string{"hamm@toy.story", "jessie@toy.story"}}
		rex    = &Participant{EmailAddress: "rex@toy.story"}
		hamm   = &Participant{EmailAddress: "hamm@toy.story", LatestRecipients: []string{"buzz@toy.story", "woody@toy.story"}}
		jessie = &Participant{EmailAddress: "jessie@toy.story", LatestRecipients: []string{"buzz@toy.story"}}

		participants = []*Participant{woody, buzz, rex, hamm, jessie}
		pairs        = []*Pair{
			&Pair{Giver: woody, Receiver: buzz},
			&Pair{Giver: buzz, Receiver: rex},
			&Pair{Giver: rex, Receiver: woody},
		}
	)

	UpdateLatestRecipients(participants, pairs)
	if len(participants) != 5 {
		t.Errorf("expected 5 participants - got %d", len(participants))
	}
	if len(woody.LatestRecipients) != 2 {
		t.Errorf("expected 2 recipients - got %d", len(woody.LatestRecipients))
	}
	if len(buzz.LatestRecipients) != 3 {
		t.Errorf("expected 3 recipients - got %d", len(buzz.LatestRecipients))
	}
	if len(rex.LatestRecipients) != 1 {
		t.Errorf("expected 1 recipient - got %d", len(rex.LatestRecipients))
	}
	if len(hamm.LatestRecipients) != 2 {
		t.Errorf("expected 2 recipients - got %d", len(hamm.LatestRecipients))
	}
	if len(jessie.LatestRecipients) != 1 {
		t.Errorf("expected 1 recipient - got %d", len(jessie.LatestRecipients))
	}
	if woody.LatestRecipients[0] != buzz.EmailAddress {
		t.Errorf("expected buzz@toy.story - got %s", woody.LatestRecipients[0])
	}
	if woody.LatestRecipients[1] != hamm.EmailAddress {
		t.Errorf("expected hamm@toy.story - got %s", woody.LatestRecipients[1])
	}
	if buzz.LatestRecipients[0] != rex.EmailAddress {
		t.Errorf("expected rex@toy.story - got %s", buzz.LatestRecipients[0])
	}
	if buzz.LatestRecipients[1] != hamm.EmailAddress {
		t.Errorf("expected hamm@toy.story - got %s", buzz.LatestRecipients[1])
	}
	if buzz.LatestRecipients[2] != jessie.EmailAddress {
		t.Errorf("expected jessie@toy.story - got %s", buzz.LatestRecipients[2])
	}
	if rex.LatestRecipients[0] != woody.EmailAddress {
		t.Errorf("expected woody@toy.story - got %s", rex.LatestRecipients[0])
	}
	if hamm.LatestRecipients[0] != buzz.EmailAddress {
		t.Errorf("expected buzz@toy.story - got %s", hamm.LatestRecipients[0])
	}
	if hamm.LatestRecipients[1] != woody.EmailAddress {
		t.Errorf("expected woody@toy.story - got %s", hamm.LatestRecipients[1])
	}
	if jessie.LatestRecipients[0] != buzz.EmailAddress {
		t.Errorf("expected buzz@toy.story - got %s", jessie.LatestRecipients[0])
	}
}

func Test_Pair_Score(t *testing.T) {
	var (
		woody  = &Participant{EmailAddress: "woody@toy.story"}
		buzz   = &Participant{EmailAddress: "buzz@toy.story"}
		rex    = &Participant{EmailAddress: "rex@toy.story", LatestRecipients: []string{"woody@toy.story"}}
		hamm   = &Participant{EmailAddress: "hamm@toy.story", LatestRecipients: []string{"not-woody@toy.story", "woody@toy.story"}}
		jessie = &Participant{EmailAddress: "jessie@toy.story", LatestRecipients: []string{"woody@toy.story", "woody@toy.story"}}
	)

	score := (&Pair{Giver: buzz, Receiver: woody}).Score()
	expectedScore := float64(0)
	if score != expectedScore {
		t.Errorf("expected score of %f - got %f", expectedScore, score)
	}

	score = (&Pair{Giver: rex, Receiver: woody}).Score()
	expectedScore = float64(1)
	if score != expectedScore {
		t.Errorf("expected score of %f - got %f", expectedScore, score)
	}

	score = (&Pair{Giver: hamm, Receiver: woody}).Score()
	expectedScore = float64(0.5)
	if score != expectedScore {
		t.Errorf("expected score of %f - got %f", expectedScore, score)
	}

	score = (&Pair{Giver: jessie, Receiver: woody}).Score()
	expectedScore = float64(1.5)
	if score != expectedScore {
		t.Errorf("expected score of %f - got %f", expectedScore, score)
	}
}

func Test_Pair_IsRepeat(t *testing.T) {
	var (
		woody  = &Participant{EmailAddress: "woody@toy.story"}
		buzz   = &Participant{EmailAddress: "buzz@toy.story"}
		rex    = &Participant{EmailAddress: "rex@toy.story", LatestRecipients: []string{"woody@toy.story"}}
		hamm   = &Participant{EmailAddress: "hamm@toy.story", LatestRecipients: []string{"not-woody@toy.story", "woody@toy.story"}}
		jessie = &Participant{EmailAddress: "jessie@toy.story", LatestRecipients: []string{"woody@toy.story", "woody@toy.story"}}
	)

	repeat := (&Pair{Giver: buzz, Receiver: woody}).IsRepeat()
	if repeat {
		t.Errorf("should not be a repeat")
	}

	repeat = (&Pair{Giver: rex, Receiver: woody}).IsRepeat()
	if !repeat {
		t.Errorf("should be a repeat")
	}

	repeat = (&Pair{Giver: hamm, Receiver: woody}).IsRepeat()
	if repeat {
		t.Errorf("should not be a repeat")
	}

	repeat = (&Pair{Giver: jessie, Receiver: woody}).IsRepeat()
	if !repeat {
		t.Errorf("should be a repeat")
	}
}
