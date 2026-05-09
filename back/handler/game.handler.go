package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/topi0247/hitori/handler/middleware"
	"github.com/topi0247/hitori/usecase"
)

type GameHandler struct {
	usecase *usecase.GameUsecase
}

func NewGameHandler(usecase *usecase.GameUsecase) *GameHandler {
	return &GameHandler{usecase: usecase}
}

func (h *GameHandler) Play(c *echo.Context) error {
	var req struct {
		ThemeID int64 `json:"theme_id"`
		Answers []struct {
			UUID  string `json:"uuid"`
			Order int    `json:"order"`
		} `json:"answers"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_request")
	}

	answers := make([]usecase.Answer, len(req.Answers))
	for i, a := range req.Answers {
		answers[i] = usecase.Answer{UUID: a.UUID, Order: a.Order}
	}

	out, err := h.usecase.Play(c.Request().Context(), usecase.PlayInput{
		ThemeID:    req.ThemeID,
		AuthUserID: middleware.AuthUserID(c),
		Answers:    answers,
	})
	if err != nil {
		log.Printf("[game.Play] error: %v", err)
		return httpError(err)
	}

	type cardResult struct {
		UUID       string `json:"uuid"`
		CardNumber int    `json:"card_number"`
		Word       string `json:"word"`
		IsCorrect  bool   `json:"is_correct"`
	}
	cards := make([]cardResult, len(out.Cards))
	for i, cr := range out.Cards {
		cards[i] = cardResult{
			UUID:       cr.UUID,
			CardNumber: cr.CardNumber,
			Word:       cr.Word,
			IsCorrect:  cr.IsCorrect,
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"correct_rate": out.CorrectRate,
		"cards":        cards,
	})
}
