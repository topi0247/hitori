package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/topi0247/hitori/usecase"
)

type ThemeHandler struct {
	usecase *usecase.ThemeUsecase
}

func NewThemeHandler(usecase *usecase.ThemeUsecase) *ThemeHandler {
	return &ThemeHandler{usecase: usecase}
}

func (h *ThemeHandler) List(c *echo.Context) error {
	out, err := h.usecase.List(c.Request().Context())
	if err != nil {
		return httpError(err)
	}

	type themeItem struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	items := make([]themeItem, len(out.Themes))
	for i, t := range out.Themes {
		items[i] = themeItem{ID: t.ID, Title: t.Title}
	}
	return c.JSON(http.StatusOK, map[string]any{"themes": items})
}
