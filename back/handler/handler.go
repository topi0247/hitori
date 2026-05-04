package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/topi0247/hitori/usecase"
)

type Handlers struct {
	Theme   *ThemeHandler
	Card    *CardHandler
	Game    *GameHandler
	Profile *ProfileHandler
}

func NewHandlers(usecases *usecase.Usecases) *Handlers {
	return &Handlers{
		Theme:   NewThemeHandler(usecases.Theme),
		Card:    NewCardHandler(usecases.Card),
		Game:    NewGameHandler(usecases.Game),
		Profile: NewProfileHandler(usecases.Profile),
	}
}

func (h *Handlers) SetRoutes(e *echo.Echo, jwtSecret string) {
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/themes", h.Theme.List)

	e.GET("/themes/:id/cards/available", h.Card.Available)
	e.GET("/themes/:id/cards/game", h.Card.GameCards)
	e.POST("/themes/:id/cards", h.Card.Create)

	e.PATCH("/cards/:id", h.Card.Confirm)
	e.DELETE("/cards/:id", h.Card.Delete)

	auth := e.Group("", JWTMiddleware(jwtSecret))
	auth.POST("/play_records", h.Game.Play)
	auth.GET("/profile", h.Profile.Get)
	auth.PATCH("/profile", h.Profile.Update)
	auth.DELETE("/profile", h.Profile.Delete)
}
