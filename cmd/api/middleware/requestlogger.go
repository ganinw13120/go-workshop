package middleware

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/pkg/log"
)

func RequestLoggerMiddleware(logger log.ILogger) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()

			var body []byte

			if c.Request().Body != nil {
				body, _ = io.ReadAll(request.Body)
			}

			logger.Info(request.Context(), "request",
				log.String("body", string(body)),
				log.String("method", request.Method),
				log.String("uri", request.RequestURI),
				log.String("remoteIP", request.RemoteAddr),
				log.String("userAgent", request.UserAgent()),
			)

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	}
}
