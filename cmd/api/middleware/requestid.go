package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/pkg/log"
)

const XRequestID string = "X-Request-Id"

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(XRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}
			newCtx := context.WithValue(c.Request().Context(), log.RequestIDKey, requestID)
			newReq := c.Request().WithContext(newCtx)
			c.SetRequest(newReq)
			c.Response().Header().Set(XRequestID, requestID)
			return next(c)
		}
	}
}
