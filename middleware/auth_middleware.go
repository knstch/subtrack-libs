package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/golang-jwt/jwt/v5"

	"github.com/knstch/subtrack-libs/auth"
	"github.com/knstch/subtrack-libs/svcerrs"
)

func WithCookieAuth(secret string) Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(http.Request)

			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				return "", svcerrs.ErrForbidden
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				return "", svcerrs.ErrForbidden
			}

			claims, err := decodeToken(secret, authHeader)
			if err != nil {
				return "", svcerrs.ErrForbidden
			}

			ctx = context.WithValue(ctx, "claims", claims)

			return next(ctx, request)
		}
	}
}

func decodeToken(secret string, token string) (*auth.Claims, error) {
	claims := &auth.Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, svcerrs.ErrForbidden
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, svcerrs.ErrUnauthorized
		}
		return claims, err
	}

	return claims, nil
}
