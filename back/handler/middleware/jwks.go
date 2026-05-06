package middleware

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
)

type jwks struct {
	Keys []struct {
		Kty string `json:"kty"`
		Alg string `json:"alg"`
		Crv string `json:"crv"`
		X   string `json:"x"`
		Y   string `json:"y"`
	} `json:"keys"`
}

func fetchECPublicKey(jwksURL string) (*ecdsa.PublicKey, error) {
	resp, err := http.Get(jwksURL) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys jwks
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}

	for _, k := range keys.Keys {
		if k.Alg == "ES256" && k.Kty == "EC" && k.Crv == "P-256" {
			xBytes, err := base64.RawURLEncoding.DecodeString(k.X)
			if err != nil {
				continue
			}
			yBytes, err := base64.RawURLEncoding.DecodeString(k.Y)
			if err != nil {
				continue
			}
			return &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     new(big.Int).SetBytes(xBytes),
				Y:     new(big.Int).SetBytes(yBytes),
			}, nil
		}
	}
	return nil, fmt.Errorf("no ES256 key found in JWKS")
}

func NewJWTConfig(hmacSecret []byte, jwksURL string) JWTConfig {
	cfg := JWTConfig{
		HMACSecret: hmacSecret,
	}
	if jwksURL != "" {
		ecKey, err := fetchECPublicKey(jwksURL)
		if err != nil {
			slog.Warn("failed to fetch JWKS", "url", jwksURL, "error", err)
			return cfg
		}
		cfg.ECPublicKey = ecKey
	}
	return cfg
}
