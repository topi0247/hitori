package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/topi0247/hitori/usecase"
)

type ProfileHandler struct {
	usecase *usecase.ProfileUsecase
}

func NewProfileHandler(usecase *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{usecase: usecase}
}

func (h *ProfileHandler) Get(c *echo.Context) error {
	out, err := h.usecase.Get(c.Request().Context(), authUserID(c))
	if err != nil {
		return httpError(err)
	}
	return c.JSON(http.StatusOK, out)
}

type updateProfileRequest struct {
	UserName string `json:"user_name"`
}

func (h *ProfileHandler) Update(c *echo.Context) error {
	var req updateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad_request")
	}

	out, err := h.usecase.Update(c.Request().Context(), usecase.UpdateProfileInput{
		AuthUserID: authUserID(c),
		UserName:   req.UserName,
	})
	if err != nil {
		return httpError(err)
	}
	return c.JSON(http.StatusOK, out)
}

func (h *ProfileHandler) Delete(c *echo.Context) error {
	if err := h.usecase.Delete(c.Request().Context(), authUserID(c)); err != nil {
		return httpError(err)
	}
	return c.NoContent(http.StatusNoContent)
}
