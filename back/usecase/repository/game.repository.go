package repository

//go:generate mockgen -source=game.repository.go -destination=mock/game_mock.go -package=repositorymock

import "context"

type PlayRecord struct {
	ThemeID     int64
	ProfileID   int64
	CardAmount  int
	CorrectRate float64
}

type PlayRecordRepository interface {
	Save(ctx context.Context, record *PlayRecord) error
}
