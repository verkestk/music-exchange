package bfscored

import (
	"fmt"
	"math/rand"
	"text/template"
	"time"

	"github.com/verkestk/music-exchange/common"
)

type pair struct {
	giver    *common.Participant
	receiver *common.Participant
}

type pairSet struct {
	pairs          []*pair
	maxScore       float64
	sumScore       float64
	minCycleLength int
}

// DoExchange matches particpants as givers and recipients, generating files with instructions for each participant
func DoExchange(participants []*common.Participant, instructionsTMPL *template.Template) error {
	pairSets := generateAllPairSets(participants)

	if len(pairSets) == 0 {
		return fmt.Errorf("no compatible receipient ordering found")
	}

	fmt.Println("number of sets", len(pairSets))
	lowestSumScorePairSets, sumScore := getLowestSumScorePairSets(pairSets)

	fmt.Println("sum score", sumScore, len(lowestSumScorePairSets))

	lowestMaxScorePairSets, maxScore := getLowestMaxScorePairSets(lowestSumScorePairSets)

	fmt.Println("max score", maxScore, len(lowestMaxScorePairSets))

	longestMinCyclePairSets, minCycleLength := getLongestMinCyclePairSets(lowestMaxScorePairSets)

	fmt.Println("min cycle length", minCycleLength, len(longestMinCyclePairSets))

	randomPairSet := getRandomPairSet(longestMinCyclePairSets)

	fmt.Printf("optimal pair set with sumScore of %f and maxScore of %f and minCycleLength of %d\n", sumScore, maxScore, minCycleLength)

	if maxScore > 0 {
		fmt.Println("A perfect pairing set was not found, so this is the best we came up with")
		fmt.Println("The following participants have had the same receipient in a past exchange:")

		for _, p := range randomPairSet.pairs {
			if p.score() > 0 {
				fmt.Println("", p.giver.ID)
			}
		}
	}

	for _, p := range randomPairSet.pairs {
		err := p.giver.WriteInstructions(p.receiver, instructionsTMPL)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateAllPairSets(participants []*common.Participant) []*pairSet {
	orders := generateAllOrders(len(participants))

	pairSets := []*pairSet{}
	for _, order := range orders {
		setPairs := []*pair{}
		isOrderValid := true
		for i, o := range order {
			if i == o {
				isOrderValid = false
				break
			}

			if !participants[i].IsCompatible(participants[o]) {
				isOrderValid = false
				break
			}
			setPairs = append(setPairs, &pair{giver: participants[i], receiver: participants[o]})
		}

		if isOrderValid {
			maxScore := float64(0)
			sumScore := float64(0)
			minCycleLength := getMinCycleLength(setPairs)

			for _, sp := range setPairs {
				score := sp.score()
				if score > maxScore {
					maxScore = score
				}
				sumScore += score
			}

			pairSets = append(pairSets, &pairSet{pairs: setPairs, maxScore: maxScore, sumScore: sumScore, minCycleLength: minCycleLength})
		}
	}

	return pairSets
}

func getLongestMinCyclePairSets(pairSets []*pairSet) (sets []*pairSet, minCycleLength int) {
	sets = []*pairSet{}
	minCycleLength = -1

	for _, ps := range pairSets {
		if minCycleLength < 0 || ps.minCycleLength > minCycleLength {
			// new lowest sumScore
			minCycleLength = ps.minCycleLength
			sets = []*pairSet{ps}
		} else if minCycleLength == ps.minCycleLength {
			sets = append(sets, ps)
		}
	}

	return sets, minCycleLength
}

func getLowestSumScorePairSets(pairSets []*pairSet) (sets []*pairSet, sumScore float64) {
	sets = []*pairSet{}
	sumScore = -1

	for _, ps := range pairSets {
		if sumScore < 0 || ps.sumScore < sumScore {
			// new lowest sumScore
			sumScore = ps.sumScore
			sets = []*pairSet{ps}
		} else if sumScore == ps.sumScore {
			sets = append(sets, ps)
		}
	}

	return sets, sumScore
}

func getLowestMaxScorePairSets(pairSets []*pairSet) (sets []*pairSet, maxScore float64) {
	sets = []*pairSet{}
	maxScore = -1

	for _, ps := range pairSets {
		if maxScore < 0 || ps.maxScore < maxScore {
			// new lowest sumScore
			maxScore = ps.maxScore
			sets = []*pairSet{ps}
		} else if maxScore == ps.maxScore {
			sets = append(sets, ps)
		}
	}

	return sets, maxScore
}

func getMinCycleLength(pairs []*pair) (minCycleLength int) {
	minCycleLength = -1

	giverIDtoPair := map[string]*pair{}
	for _, p := range pairs {
		giverIDtoPair[p.giver.ID] = p
	}

	for _, p := range pairs {
		cycleLength := 1
		originalGiver := p.giver
		currentPair := p
		for originalGiver.ID != currentPair.receiver.ID {
			cycleLength++
			currentPair = giverIDtoPair[currentPair.receiver.ID]
		}

		if minCycleLength == -1 || cycleLength < minCycleLength {
			minCycleLength = cycleLength
		}
	}

	return minCycleLength
}

func getRandomPairSet(pairSets []*pairSet) *pairSet {
	rand.Seed(time.Now().UnixNano())
	return pairSets[rand.Intn(len(pairSets))]
}

func (p *pair) score() float64 {
	score := float64(0)
	for i, previousRecipient := range p.giver.LatestRecipients {
		if previousRecipient == p.receiver.ID {
			score += float64(1) / float64(i+1)
		}
	}

	return score
}

func generateAllOrders(length int) [][]int {
	inOrder := []int{}
	for i := 0; i < length; i++ {
		inOrder = append(inOrder, i)
	}

	return permutations(inOrder)
}

// from https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
