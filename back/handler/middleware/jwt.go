package middleware

import (
	"crypto/ecdsa"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

const authUserIDKey = "auth_user_id"

type JWTConfig struct {
	HMACSecret  []byte
	ECPublicKey *ecdsa.PublicKey
}

func (cfg *JWTConfig) keyfunc(t *jwt.Token) (any, error) {
	switch t.Method.(type) {
	case *jwt.SigningMethodHMAC:
		return cfg.HMACSecret, nil
	case *jwt.SigningMethodECDSA:
		if cfg.ECPublicKey == nil {
			return nil, jwt.ErrSignatureInvalid
		}
		return cfg.ECPublicKey, nil
	default:
		return nil, jwt.ErrSignatureInvalid
	}
}

func JWT(cfg JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, cfg.keyfunc)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			sub, ok := claims["sub"].(string)
			if !ok || sub == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			c.Set(authUserIDKey, sub)
			return next(c)
		}
	}
}

func OptionalJWT(cfg JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
				token, err := jwt.Parse(tokenStr, cfg.keyfunc)
				if err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						if sub, ok := claims["sub"].(string); ok && sub != "" {
							c.Set(authUserIDKey, sub)
						}
					}
				}
			}
			return next(c)
		}
	}
}

func AuthUserID(c *echo.Context) string {
	v, _ := c.Get(authUserIDKey).(string)
	return v
}
