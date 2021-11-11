package bfscored

import (
	"testing"

	"github.com/verkestk/music-exchange/src/participant"
)

func Test_DoExchange(t *testing.T) {
	// participants []*participant.Participant, allowRepeatRecipients bool
	woody := &participant.Participant{EmailAddress: "woody", Platforms: []string{"TheClaw"}, LatestRecipients: []string{"buzz"}}
	buzz := &participant.Participant{EmailAddress: "buzz", Platforms: []string{"TheClaw"}, LatestRecipients: []string{"woody"}}
	rex := &participant.Participant{EmailAddress: "rex", Platforms: []string{"DayCare"}}
	hamm := &participant.Participant{EmailAddress: "hamm", Platforms: []string{"DayCare"}}
	jessie := &participant.Participant{EmailAddress: "jessie", Platforms: []string{"DayCare"}}

	participants := []*participant.Participant{woody, buzz, rex, hamm, jessie}

	pairs, err := DoExchange(participants, true)
	if err != nil {
		t.Errorf("unexpected error: %w", err)
	}
	if len(pairs) != 5 {
		t.Errorf("expected %d pairs - got %d", 5, len(pairs))
	}

	_, err = DoExchange(participants, false)
	if err == nil {
		t.Errorf("expected error: %w", err)
	}
}

func Test_getLongestMinCyclePairSets(t *testing.T) {
	woody := &participant.Participant{EmailAddress: "woody", Platforms: []string{"TheClaw"}}
	buzz := &participant.Participant{EmailAddress: "buzz", Platforms: []string{"TheClaw"}}
	rex := &participant.Participant{EmailAddress: "rex", Platforms: []string{"TheClaw"}}
	hamm := &participant.Participant{EmailAddress: "hamm", Platforms: []string{"TheClaw"}}
	jessie := &participant.Participant{EmailAddress: "jessie", Platforms: []string{"TheClaw"}}

	participants := []*participant.Participant{
		woody,
		buzz,
		rex,
		hamm,
		jessie,
	}

	pairSets := generateAllPairSets(participants, true)

	longestMinCyclePairSets, minCycleLength := getLongestMinCyclePairSets(pairSets)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}
	if len(longestMinCyclePairSets) != 24 {
		t.Errorf("len(longestMinCyclePairSets) = %d; want 24", len(longestMinCyclePairSets))
	}
}

func Test_getMinCycleLength(t *testing.T) {
	woody := &participant.Participant{EmailAddress: "woody"}
	buzz := &participant.Participant{EmailAddress: "buzz"}
	rex := &participant.Participant{EmailAddress: "rex"}
	hamm := &participant.Participant{EmailAddress: "hamm"}
	jessie := &participant.Participant{EmailAddress: "jessie"}

	pairs := []*participant.Pair{
		&participant.Pair{Giver: woody, Receiver: rex},
		&participant.Pair{Giver: rex, Receiver: hamm},
		&participant.Pair{Giver: hamm, Receiver: woody},
		&participant.Pair{Giver: buzz, Receiver: jessie},
		&participant.Pair{Giver: jessie, Receiver: buzz},
	}

	minCycleLength := getMinCycleLength(pairs)
	if minCycleLength != 2 {
		t.Errorf("minCycleLength = %d; want 2", minCycleLength)
	}

	// same pairs, different order in the slice
	pairs = []*participant.Pair{
		&participant.Pair{Giver: woody, Receiver: rex},
		&participant.Pair{Giver: buzz, Receiver: jessie},
		&participant.Pair{Giver: rex, Receiver: hamm},
		&participant.Pair{Giver: hamm, Receiver: woody},
		&participant.Pair{Giver: jessie, Receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 2 {
		t.Errorf("minCycleLength = %d; want 2", minCycleLength)
	}

	// long cycle
	pairs = []*participant.Pair{
		&participant.Pair{Giver: woody, Receiver: buzz},
		&participant.Pair{Giver: buzz, Receiver: rex},
		&participant.Pair{Giver: rex, Receiver: hamm},
		&participant.Pair{Giver: hamm, Receiver: jessie},
		&participant.Pair{Giver: jessie, Receiver: woody},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}

	// long cycle, different order
	pairs = []*participant.Pair{
		&participant.Pair{Giver: buzz, Receiver: rex},
		&participant.Pair{Giver: hamm, Receiver: jessie},
		&participant.Pair{Giver: jessie, Receiver: woody},
		&participant.Pair{Giver: rex, Receiver: hamm},
		&participant.Pair{Giver: woody, Receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}

	// cycle of 1?
	pairs = []*participant.Pair{
		&participant.Pair{Giver: buzz, Receiver: rex},
		&participant.Pair{Giver: hamm, Receiver: woody},
		&participant.Pair{Giver: jessie, Receiver: jessie},
		&participant.Pair{Giver: rex, Receiver: hamm},
		&participant.Pair{Giver: woody, Receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 1 {
		t.Errorf("minCycleLength = %d; want 1", minCycleLength)
	}
}
