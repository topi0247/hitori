package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	domainCard "github.com/topi0247/hitori/domain/card"
	domainProfile "github.com/topi0247/hitori/domain/profile"
	domainTheme "github.com/topi0247/hitori/domain/theme"
	"github.com/topi0247/hitori/usecase"
)

func httpError(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, domainTheme.ErrNotFound),
		errors.Is(err, domainCard.ErrNotFound),
		errors.Is(err, domainProfile.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "not_found")

	case errors.Is(err, domainCard.ErrInvalidCardNumber),
		errors.Is(err, domainCard.ErrInvalidWord),
		errors.Is(err, domainCard.ErrInvalidGuestName),
		errors.Is(err, domainCard.ErrOwnerRequired),
		errors.Is(err, domainProfile.ErrInvalidUserName),
		errors.Is(err, usecase.ErrInvalidCardAmount),
		errors.Is(err, usecase.ErrNoAnswers):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	case errors.Is(err, domainCard.ErrForbidden):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())

	case errors.Is(err, domainCard.ErrAlreadyConfirmed),
		errors.Is(err, domainCard.ErrThemeCardLimitReached):
		return echo.NewHTTPError(http.StatusConflict, err.Error())

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "unexpected_error")
	}
}
