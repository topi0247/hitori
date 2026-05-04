package usecase

import "github.com/topi0247/hitori/usecase/repository"

type Usecases struct {
	Card  *CardUsecase
	Theme *ThemeUsecase
	Game  *GameUsecase
}

func NewUsecases(repos *repository.Repositories) *Usecases {
	return &Usecases{
		Card:  NewCardUsecase(repos.Card, repos.Theme),
		Theme: NewThemeUsecase(repos.Theme),
		Game:  NewGameUsecase(repos.Card, repos.PlayRecord),
	}
}
