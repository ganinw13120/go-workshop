package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/config"

	"github.com/wisesight/go-api-template/cmd/api/handler"
	"github.com/wisesight/go-api-template/cmd/api/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/wisesight/go-api-template/cmd/api/docs" // docs is generated by Swag CLI, you have to import it.
)

func NewRoute(
	config config.Config,
	app *echo.Echo,
	probeHandler handler.IProbe,
	timelineHandler handler.ITimeline,
	accountHandler handler.IAccount,
) {
	app.GET("/", func(c echo.Context) error {

		return c.String(http.StatusOK, "Hello world")
	})
	app.GET("/readyz", probeHandler.DBReadyCheck)

	u := app.Group("/user")

	u.Use(middleware.NewVerifyJWTAuth([]byte(config.JWTSecret), config.JWTSigningMethod))
	u.Use(middleware.ExtractJWTClaims)

	t := app.Group("/timeline")
	t.POST("/", timelineHandler.Save)
	t.GET("/", timelineHandler.Get)

	account := app.Group("/account")
	account.PUT("/", accountHandler.Save)

	app.GET("/swagger/*", echoSwagger.WrapHandler)
}
