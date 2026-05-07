package usecase

import (
	"context"
	"errors"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/usecase/repository"
)

var ErrNoAnswers = errors.New("no_answers")

type GameUsecase struct {
	cardRepository       repository.CardRepository
	playRecordRepository repository.PlayRecordRepository
	profileRepository    repository.ProfileRepository
}

func NewGameUsecase(cardRepository repository.CardRepository, playRecordRepository repository.PlayRecordRepository, profileRepository repository.ProfileRepository) *GameUsecase {
	return &GameUsecase{cardRepository: cardRepository, playRecordRepository: playRecordRepository, profileRepository: profileRepository}
}

// --- Play ---

type PlayInput struct {
	ThemeID    int64
	AuthUserID string
	Answers    []Answer
}

type Answer struct {
	UUID  string
	Order int
}

type PlayCardResult struct {
	UUID       string
	CardNumber int
	Word       string
	IsCorrect  bool
}

type PlayOutput struct {
	CorrectRate float64
	Cards       []PlayCardResult
}

func (u *GameUsecase) Play(ctx context.Context, input PlayInput) (*PlayOutput, error) {
	if len(input.Answers) == 0 {
		return nil, ErrNoAnswers
	}

	profile, err := u.profileRepository.FetchByAuthUserID(ctx, input.AuthUserID)
	if err != nil {
		return nil, err
	}

	uuids := make([]string, len(input.Answers))
	for i, a := range input.Answers {
		uuids[i] = a.UUID
	}

	cards, err := u.cardRepository.FetchGameCardsByUUIDs(ctx, uuids)
	if err != nil {
		return nil, err
	}

	entries := make([]domain.CardEntry, len(cards))
	for i, c := range cards {
		entries[i] = domain.CardEntry{UUID: c.UUID, CardNumber: c.Number}
	}

	domainAnswers := make([]domain.GameAnswer, len(input.Answers))
	for i, a := range input.Answers {
		domainAnswers[i] = domain.GameAnswer{UUID: a.UUID, Order: a.Order}
	}

	result, err := domain.Judge(entries, domainAnswers)
	if err != nil {
		return nil, err
	}

	wordMap := make(map[string]string, len(cards))
	for _, c := range cards {
		wordMap[c.UUID] = c.Word
	}

	for _, cr := range result.Cards {
		if cr.MatchPoints > 0 {
			if err := u.cardRepository.AddMatchPoints(ctx, cr.UUID, cr.MatchPoints); err != nil {
				return nil, err
			}
		}
	}

	if err := u.playRecordRepository.Save(ctx, &repository.PlayRecord{
		ThemeID:     input.ThemeID,
		ProfileID:   profile.ID,
		CardAmount:  len(cards),
		CorrectRate: result.CorrectRate,
	}); err != nil {
		return nil, err
	}

	cardResults := make([]PlayCardResult, len(result.Cards))
	for i, cr := range result.Cards {
		cardResults[i] = PlayCardResult{
			UUID:       cr.UUID,
			CardNumber: cr.CardNumber,
			Word:       wordMap[cr.UUID],
			IsCorrect:  cr.IsCorrect,
		}
	}

	return &PlayOutput{CorrectRate: result.CorrectRate, Cards: cardResults}, nil
}
