package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/topi0247/hitori/handler/middleware"
	"github.com/topi0247/hitori/usecase"
)

type CardHandler struct {
	usecase *usecase.CardUsecase
}

func NewCardHandler(usecase *usecase.CardUsecase) *CardHandler {
	return &CardHandler{usecase: usecase}
}

func (h *CardHandler) Available(c *echo.Context) error {
	themeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_theme_id")
	}

	out, err := h.usecase.Available(c.Request().Context(), themeID)
	if err != nil {
		return httpError(err)
	}
	return c.JSON(http.StatusOK, map[string]any{"card_number": out.CardNumber})
}

func (h *CardHandler) GameCards(c *echo.Context) error {
	themeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_theme_id")
	}
	cardAmount, err := strconv.Atoi(c.QueryParam("card_amount"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_card_amount")
	}

	out, err := h.usecase.GameCards(c.Request().Context(), usecase.GameCardsInput{
		ThemeID:    themeID,
		CardAmount: cardAmount,
	})
	if err != nil {
		return httpError(err)
	}

	type cardItem struct {
		UUID string `json:"uuid"`
		Word string `json:"word"`
	}
	items := make([]cardItem, len(out.Cards))
	for i, card := range out.Cards {
		items[i] = cardItem{UUID: card.UUID, Word: card.Word}
	}
	return c.JSON(http.StatusOK, map[string]any{"cards": items})
}

func (h *CardHandler) Create(c *echo.Context) error {
	themeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_theme_id")
	}

	var req struct {
		CardNumber int     `json:"card_number"`
		Word       string  `json:"word"`
		GuestName  *string `json:"guest_name"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_request")
	}

	out, err := h.usecase.Create(c.Request().Context(), usecase.CreateInput{
		ThemeID:    themeID,
		CardNumber: req.CardNumber,
		Word:       req.Word,
		AuthUserID: middleware.AuthUserID(c),
		GuestName:  req.GuestName,
		Now:        time.Now(),
	})
	if err != nil {
		return httpError(err)
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"id":          out.ID,
		"uuid":        out.UUID,
		"card_number": out.CardNumber,
		"word":        out.Word,
	})
}

func (h *CardHandler) Confirm(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_id")
	}

	var req struct {
		Word      string  `json:"word"`
		GuestName *string `json:"guest_name"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_request")
	}

	out, err := h.usecase.Confirm(c.Request().Context(), usecase.ConfirmInput{
		ID:         id,
		Word:       req.Word,
		AuthUserID: middleware.AuthUserID(c),
		GuestName:  req.GuestName,
	})
	if err != nil {
		return httpError(err)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"id":          out.ID,
		"card_number": out.CardNumber,
		"word":        out.Word,
	})
}

func (h *CardHandler) Delete(c *echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_id")
	}

	var req struct {
		GuestName *string `json:"guest_name"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid_request")
	}

	if err := h.usecase.Delete(c.Request().Context(), usecase.DeleteInput{
		ID:         id,
		AuthUserID: middleware.AuthUserID(c),
		GuestName:  req.GuestName,
	}); err != nil {
		return httpError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
