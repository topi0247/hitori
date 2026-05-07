package domain

import (
	"errors"
	"math"
	"sort"
)

const (
	PointsExact    = 3
	PointsAdjacent = 2
)

type CardEntry struct {
	UUID       string
	CardNumber int
}

type GameAnswer struct {
	UUID  string
	Order int
}

type CardResult struct {
	UUID        string
	CardNumber  int
	IsCorrect   bool
	MatchPoints int
}

type JudgeResult struct {
	CorrectRate float64
	Cards       []CardResult
}

func Judge(cards []CardEntry, answers []GameAnswer) (*JudgeResult, error) {
	cardMap := make(map[string]CardEntry, len(cards))
	for _, c := range cards {
		cardMap[c.UUID] = c
	}
	for _, a := range answers {
		if _, ok := cardMap[a.UUID]; !ok {
			return nil, errors.New("unknown_uuid")
		}
	}

	sorted := make([]CardEntry, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CardNumber < sorted[j].CardNumber
	})
	correctOrder := make(map[string]int, len(sorted))
	for i, c := range sorted {
		correctOrder[c.UUID] = i + 1
	}

	userOrder := make(map[string]int, len(answers))
	for _, a := range answers {
		userOrder[a.UUID] = a.Order
	}

	correctCount := 0
	results := make([]CardResult, 0, len(cards))

	for _, c := range cards {
		correct := correctOrder[c.UUID]
		user := userOrder[c.UUID]
		diff := int(math.Abs(float64(correct - user)))

		var pts int
		var isCorrect bool
		switch diff {
		case 0:
			pts = PointsExact
			isCorrect = true
			correctCount++
		case 1:
			pts = PointsAdjacent
		}

		results = append(results, CardResult{
			UUID:        c.UUID,
			CardNumber:  c.CardNumber,
			IsCorrect:   isCorrect,
			MatchPoints: pts,
		})
	}

	correctRate := math.Round(float64(correctCount)/float64(len(cards))*100*100) / 100

	return &JudgeResult{
		CorrectRate: correctRate,
		Cards:       results,
	}, nil
}
