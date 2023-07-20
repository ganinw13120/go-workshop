package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wisesight/go-api-template/pkg/log"
)

func ResponseLoggerMiddleware(logger log.ILogger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogError:        true,
		LogMethod:       true,
		LogRemoteIP:     true,
		LogUserAgent:    true,
		LogResponseSize: true,
		Skipper:         middleware.DefaultSkipper,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			request := c.Request()
			logger.Info(request.Context(), "request",
				log.String("method", v.Method),
				log.String("uri", v.URI),
				log.Int("status", v.Status),
				log.String("remoteIP", v.RemoteIP),
				log.String("userAgent", v.UserAgent),
				log.Int64("bytesOut", v.ResponseSize),
			)

			return nil
		},
	})
}
