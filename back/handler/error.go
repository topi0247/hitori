package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/topi0247/hitori/domain"
	"github.com/topi0247/hitori/usecase"
)

func httpError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "not_found")

	case errors.Is(err, domain.ErrInvalidCardNumber),
		errors.Is(err, domain.ErrInvalidWord),
		errors.Is(err, domain.ErrInvalidGuestName),
		errors.Is(err, domain.ErrOwnerRequired),
		errors.Is(err, domain.ErrInvalidUserName),
		errors.Is(err, usecase.ErrInvalidCardAmount),
		errors.Is(err, usecase.ErrNoAnswers):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	case errors.Is(err, domain.ErrForbidden):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())

	case errors.Is(err, domain.ErrAlreadyConfirmed),
		errors.Is(err, domain.ErrThemeCardLimitReached),
		errors.Is(err, domain.ErrAlreadyExists):
		return echo.NewHTTPError(http.StatusConflict, err.Error())

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected_error")
	}
}
