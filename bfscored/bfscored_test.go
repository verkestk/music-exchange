package bfscored

import (
	"testing"

	"github.com/verkestk/music-exchange/common"
)

func Test_getLongestMinCyclePairSets(t *testing.T) {
	woody := &common.Participant{ID: "woody", Platforms: []string{"TheClaw"}}
	buzz := &common.Participant{ID: "buzz", Platforms: []string{"TheClaw"}}
	rex := &common.Participant{ID: "rex", Platforms: []string{"TheClaw"}}
	hamm := &common.Participant{ID: "hamm", Platforms: []string{"TheClaw"}}
	jessie := &common.Participant{ID: "jessie", Platforms: []string{"TheClaw"}}

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

func Test_getMinCycleLength(t *testing.T) {
	woody := &common.Participant{ID: "woody"}
	buzz := &common.Participant{ID: "buzz"}
	rex := &common.Participant{ID: "rex"}
	hamm := &common.Participant{ID: "hamm"}
	jessie := &common.Participant{ID: "jessie"}

	pairs := []*pair{
		&pair{giver: woody, receiver: rex},
		&pair{giver: rex, receiver: hamm},
		&pair{giver: hamm, receiver: woody},
		&pair{giver: buzz, receiver: jessie},
		&pair{giver: jessie, receiver: buzz},
	}

	minCycleLength := getMinCycleLength(pairs)
	if minCycleLength != 2 {
		t.Errorf("minCycleLength = %d; want 2", minCycleLength)
	}

	// same pairs, different order in the slice
	pairs = []*pair{
		&pair{giver: woody, receiver: rex},
		&pair{giver: buzz, receiver: jessie},
		&pair{giver: rex, receiver: hamm},
		&pair{giver: hamm, receiver: woody},
		&pair{giver: jessie, receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 2 {
		t.Errorf("minCycleLength = %d; want 2", minCycleLength)
	}

	// long cycle
	pairs = []*pair{
		&pair{giver: woody, receiver: buzz},
		&pair{giver: buzz, receiver: rex},
		&pair{giver: rex, receiver: hamm},
		&pair{giver: hamm, receiver: jessie},
		&pair{giver: jessie, receiver: woody},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}

	// long cycle, different order
	pairs = []*pair{
		&pair{giver: buzz, receiver: rex},
		&pair{giver: hamm, receiver: jessie},
		&pair{giver: jessie, receiver: woody},
		&pair{giver: rex, receiver: hamm},
		&pair{giver: woody, receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 5 {
		t.Errorf("minCycleLength = %d; want 5", minCycleLength)
	}

	// cycle of 1?
	pairs = []*pair{
		&pair{giver: buzz, receiver: rex},
		&pair{giver: hamm, receiver: woody},
		&pair{giver: jessie, receiver: jessie},
		&pair{giver: rex, receiver: hamm},
		&pair{giver: woody, receiver: buzz},
	}

	minCycleLength = getMinCycleLength(pairs)
	if minCycleLength != 1 {
		t.Errorf("minCycleLength = %d; want 1", minCycleLength)
	}
}
