package usecase

import (
	"context"
	"errors"
	"time"

	domainCard "github.com/topi0247/hitori/domain/card"
	"github.com/topi0247/hitori/usecase/repository"
)

const (
	MinCardAmount = 4
	MaxCardAmount = 10
)

var ErrInvalidCardAmount = errors.New("invalid_card_amount")

type CardUsecase struct {
	cardRepository    repository.CardRepository
	themeRepository   repository.ThemeRepository
	profileRepository repository.ProfileRepository
}

func NewCardUsecase(cardRepository repository.CardRepository, themeRepository repository.ThemeRepository, profileRepository repository.ProfileRepository) *CardUsecase {
	return &CardUsecase{
		cardRepository:    cardRepository,
		themeRepository:   themeRepository,
		profileRepository: profileRepository,
	}
}

func (u *CardUsecase) resolveProfileID(ctx context.Context, authUserID string) (*int64, error) {
	if authUserID == "" {
		return nil, nil
	}
	p, err := u.profileRepository.FetchByAuthUserID(ctx, authUserID)
	if err != nil {
		return nil, err
	}
	return &p.ID, nil
}

// --- Create ---

type CreateInput struct {
	ThemeID    int64
	CardNumber int
	Word       string
	AuthUserID string
	GuestName  *string
	Now        time.Time
}

type CreateOutput struct {
	ID         int64
	CardNumber int
	Word       string
}

func (u *CardUsecase) Create(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if _, err := u.themeRepository.FetchByID(ctx, input.ThemeID); err != nil {
		return nil, err
	}

	count, err := u.cardRepository.CountByThemeID(ctx, input.ThemeID)
	if err != nil {
		return nil, err
	}
	if err := domainCard.CanAddCard(count); err != nil {
		return nil, err
	}

	profileID, err := u.resolveProfileID(ctx, input.AuthUserID)
	if err != nil {
		return nil, err
	}

	c, err := domainCard.New(domainCard.NewInput{
		ThemeID:   input.ThemeID,
		Number:    input.CardNumber,
		Word:      input.Word,
		ProfileID: profileID,
		GuestName: input.GuestName,
		Now:       input.Now,
	})
	if err != nil {
		return nil, err
	}

	if err := u.cardRepository.Save(ctx, c); err != nil {
		return nil, err
	}

	return &CreateOutput{
		ID:         c.ID,
		CardNumber: c.Number,
		Word:       c.Word,
	}, nil
}

// --- Confirm ---

type ConfirmInput struct {
	ID         int64
	Word       string
	AuthUserID string
	GuestName  *string
}

type ConfirmOutput struct {
	ID         int64
	CardNumber int
	Word       string
}

func (u *CardUsecase) Confirm(ctx context.Context, input ConfirmInput) (*ConfirmOutput, error) {
	c, err := u.cardRepository.FetchByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	profileID, err := u.resolveProfileID(ctx, input.AuthUserID)
	if err != nil {
		return nil, err
	}

	if err := c.VerifyOwner(profileID, input.GuestName); err != nil {
		return nil, err
	}

	if err := c.Confirm(input.Word); err != nil {
		return nil, err
	}

	if err := u.cardRepository.Confirm(ctx, input.ID, input.Word); err != nil {
		return nil, err
	}

	return &ConfirmOutput{
		ID:         c.ID,
		CardNumber: c.Number,
		Word:       c.Word,
	}, nil
}

// --- Delete ---

type DeleteInput struct {
	ID         int64
	AuthUserID string
	GuestName  *string
}

func (u *CardUsecase) Delete(ctx context.Context, input DeleteInput) error {
	c, err := u.cardRepository.FetchByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if c.IsConfirmed {
		return domainCard.ErrAlreadyConfirmed
	}

	profileID, err := u.resolveProfileID(ctx, input.AuthUserID)
	if err != nil {
		return err
	}

	if err := c.VerifyOwner(profileID, input.GuestName); err != nil {
		return err
	}

	return u.cardRepository.Delete(ctx, input.ID)
}

// --- Available ---

type AvailableOutput struct {
	CardNumber int
}

func (u *CardUsecase) Available(ctx context.Context, themeID int64) (*AvailableOutput, error) {
	if _, err := u.themeRepository.FetchByID(ctx, themeID); err != nil {
		return nil, err
	}
	number, err := u.cardRepository.GetAvailableNumber(ctx, themeID)
	if err != nil {
		return nil, err
	}
	return &AvailableOutput{CardNumber: number}, nil
}

// --- GameCards ---

type GameCardsInput struct {
	ThemeID    int64
	CardAmount int
}

type GameCardsOutput struct {
	Cards []*domainCard.Card
}

func (u *CardUsecase) GameCards(ctx context.Context, input GameCardsInput) (*GameCardsOutput, error) {
	if input.CardAmount < MinCardAmount || input.CardAmount > MaxCardAmount {
		return nil, ErrInvalidCardAmount
	}
	if _, err := u.themeRepository.FetchByID(ctx, input.ThemeID); err != nil {
		return nil, err
	}
	cards, err := u.cardRepository.GetGameCards(ctx, input.ThemeID, input.CardAmount)
	if err != nil {
		return nil, err
	}
	return &GameCardsOutput{Cards: cards}, nil
}
