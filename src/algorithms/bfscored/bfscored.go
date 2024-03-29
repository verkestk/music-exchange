package bfscored

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/verkestk/music-exchange/src/participant"
)

type pairSet struct {
	pairs          []*participant.Pair
	maxScore       float64
	sumScore       float64
	minCycleLength int
}

// DoExchange matches particpants as givers and recipients, generating files with instructions for each participant
func DoExchange(participants []*participant.Participant, allowRepeatRecipients bool) ([]*participant.Pair, error) {
	pairSets := generateAllPairSets(participants, allowRepeatRecipients)

	if len(pairSets) == 0 {
		return nil, fmt.Errorf("no compatible receipient ordering found for %d participants", len(participants))
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
			if p.Score() > 0 {
				fmt.Println("", p.Giver.EmailAddress)
			}
		}
	}

	return randomPairSet.pairs, nil
}

func generateAllPairSets(participants []*participant.Participant, allowRepeatRecipients bool) []*pairSet {
	orders := generateAllOrders(len(participants))

	pairSets := []*pairSet{}
	for _, order := range orders {
		setPairs := []*participant.Pair{}
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
			setPairs = append(setPairs, &participant.Pair{Giver: participants[i], Receiver: participants[o]})
		}

		if isOrderValid {
			maxScore := float64(0)
			sumScore := float64(0)
			minCycleLength := getMinCycleLength(setPairs)

			hasRepeat := false

			for _, sp := range setPairs {
				score := sp.Score()

				hasRepeat = hasRepeat || sp.IsRepeat()

				if score > maxScore {
					maxScore = score
				}
				sumScore += score
			}

			if !hasRepeat || allowRepeatRecipients {
				pairSets = append(pairSets, &pairSet{pairs: setPairs, maxScore: maxScore, sumScore: sumScore, minCycleLength: minCycleLength})
			}
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

func getMinCycleLength(pairs []*participant.Pair) (minCycleLength int) {
	minCycleLength = -1

	giverEmailAddresstoPair := map[string]*participant.Pair{}
	for _, p := range pairs {
		giverEmailAddresstoPair[p.Giver.EmailAddress] = p
	}

	for _, p := range pairs {
		cycleLength := 1
		originalGiver := p.Giver
		currentPair := p
		for originalGiver.EmailAddress != currentPair.Receiver.EmailAddress {
			cycleLength++
			currentPair = giverEmailAddresstoPair[currentPair.Receiver.EmailAddress]
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
