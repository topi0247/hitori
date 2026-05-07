package usecase

import "github.com/topi0247/hitori/domain"

type Usecases struct {
	Card    *CardUsecase
	Theme   *ThemeUsecase
	Game    *GameUsecase
	Profile *ProfileUsecase
}

func NewUsecases(repos *domain.Repositories) *Usecases {
	return &Usecases{
		Card:    NewCardUsecase(repos.Card, repos.Theme, repos.Profile),
		Theme:   NewThemeUsecase(repos.Theme),
		Game:    NewGameUsecase(repos.Card, repos.PlayRecord, repos.Profile),
		Profile: NewProfileUsecase(repos.Profile),
	}
}
