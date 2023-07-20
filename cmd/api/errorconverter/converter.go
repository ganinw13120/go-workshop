package errorconverter

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/wisesight/go-api-template/pkg/apperror"
)

func ResponseError(c echo.Context, err error) error {
	type ErrorResponse struct {
		Message string
	}

	switch e := err.(type) {
	case apperror.AppError:
		// Code แต่ละ handler จะมีความหมายไม่เหมือนกัน อาจจะใช้ coverter เดียวกันไม่ได้
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: e.Message,
		})
	case error:
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Internal server error",
		})
	}
	return nil
}
