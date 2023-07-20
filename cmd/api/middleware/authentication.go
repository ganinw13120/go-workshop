package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/wisesight/go-api-template/constant"
	"github.com/wisesight/go-api-template/pkg/entity"
)

func NewVerifyJWTAuth(secret interface{}, signingMethod string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    secret,
		SigningMethod: signingMethod,
		ContextKey:    constant.JWT_CONTEXT_KEY,
	})
}

func ExtractJWTClaims(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get(constant.JWT_CONTEXT_KEY).(*jwt.Token)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.NoContent(http.StatusUnauthorized)
		}

		var user entity.UserSession
		mapstructure.Decode(claims, &user)

		c.Set(constant.JWT_CONTEXT_KEY, user)
		return next(c)
	}
}
