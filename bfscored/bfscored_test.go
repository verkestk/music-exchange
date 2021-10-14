package bfscored

import (
	"testing"

	"github.com/verkestk/music-exchange/common"
)

func Test_generateAllPairSets(t *testing.T) {
	// TODO
}

func Test_getLongestMinCyclePairSets(t *testing.T) {
	woody := &common.Participant{EmailAddress: "woody", Platforms: []string{"TheClaw"}}
	buzz := &common.Participant{EmailAddress: "buzz", Platforms: []string{"TheClaw"}}
	rex := &common.Participant{EmailAddress: "rex", Platforms: []string{"TheClaw"}}
	hamm := &common.Participant{EmailAddress: "hamm", Platforms: []string{"TheClaw"}}
	jessie := &common.Participant{EmailAddress: "jessie", Platforms: []string{"TheClaw"}}

	participants := []*common.Participant{
		woody,
		buzz,
		rex,
		hamm,
		jessie,
	}

	pairSets := generateAllPairSets(participants)

	longestMinCyclePairSets, minCycleLength := getLongestMinCyclePairSets(pairSets)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}
	if len(longestMinCyclePairSets) != 24 {
		t.Errorf("len(longestMinCyclePairSets) = %d; want 24", len(longestMinCyclePairSets))
	}
}

func Test_getLowestSumScorePairSets(t *testing.T) {
	// TODO
}

func Test_getLowestMaxScorePairSets(t *testing.T) {
	// TODO
}

func Test_getMinCycleLength(t *testing.T) {
	woody := &common.Participant{EmailAddress: "woody"}
	buzz := &common.Participant{EmailAddress: "buzz"}
	rex := &common.Participant{EmailAddress: "rex"}
	hamm := &common.Participant{EmailAddress: "hamm"}
	jessie := &common.Participant{EmailAddress: "jessie"}

	pairs := []*common.Pair{
		&common.Pair{Giver: woody, Receiver: rex},
		&common.Pair{Giver: rex, Receiver: hamm},
		&common.Pair{Giver: hamm, Receiver: woody},
		&common.Pair{Giver: buzz, Receiver: jessie},
		&common.Pair{Giver: jessie, Receiver: buzz},
	}

	minCycleLength := getMinCycleLength(pairs)
	if minCycleLength != 2 {
		t.Errorf("minCycleLength = %d; want 2", minCycleLength)
	}

	// same pairs, different order in the slice
	pairs = []*common.Pair{
		&common.Pair{Giver: woody, Receiver: rex},
		&common.Pair{Giver: buzz, Receiver: jessie},
		&common.Pair{Giver: rex, Receiver: hamm},
		&common.Pair{Giver: hamm, Receiver: woody},
		&common.Pair{Giver: jessie, Receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 2 {
		t.Errorf("minCycleLength = %d; want 2", minCycleLength)
	}

	// long cycle
	pairs = []*common.Pair{
		&common.Pair{Giver: woody, Receiver: buzz},
		&common.Pair{Giver: buzz, Receiver: rex},
		&common.Pair{Giver: rex, Receiver: hamm},
		&common.Pair{Giver: hamm, Receiver: jessie},
		&common.Pair{Giver: jessie, Receiver: woody},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}

	// long cycle, different order
	pairs = []*common.Pair{
		&common.Pair{Giver: buzz, Receiver: rex},
		&common.Pair{Giver: hamm, Receiver: jessie},
		&common.Pair{Giver: jessie, Receiver: woody},
		&common.Pair{Giver: rex, Receiver: hamm},
		&common.Pair{Giver: woody, Receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}

	// cycle of 1?
	pairs = []*common.Pair{
		&common.Pair{Giver: buzz, Receiver: rex},
		&common.Pair{Giver: hamm, Receiver: woody},
		&common.Pair{Giver: jessie, Receiver: jessie},
		&common.Pair{Giver: rex, Receiver: hamm},
		&common.Pair{Giver: woody, Receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 1 {
		t.Errorf("minCycleLength = %d; want 1", minCycleLength)
	}
}

func Test_getRandomPairSet(t *testing.T) {
	// TODO
}

func Test_pair_score(t *testing.T) {
	// TODO
}

func Test_generateAllOrders(t *testing.T) {
	// TODO
}

func Test_permutations(t *testing.T) {
	// TODO
}
