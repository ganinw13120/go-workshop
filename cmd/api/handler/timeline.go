package handler

import (
	gpgvalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/helper"
	"github.com/wisesight/go-api-template/pkg/log"
	"github.com/wisesight/go-api-template/pkg/usecase"
	"github.com/wisesight/go-api-template/pkg/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type ITimeline interface {
	Save(c echo.Context) error
}

type timeline struct {
	timelineUseCase usecase.ITimeline
	logger          log.ILogger
}

func NewTimeline(timelineUseCase usecase.ITimeline, logger log.ILogger) *timeline {
	return &timeline{
		timelineUseCase: timelineUseCase,
		logger:          logger,
	}
}

type SaveRequestBody struct {
	Text        string `json:"name" validate:"required"`
	UserId      string `json:"user_id" validate:"required"`
	Likes       int    `json:"likes" validate:"required"`
	RepostCount int    `json:"repost_count" validate:"required"`
}

func (t timeline) Save(c echo.Context) error {
	body := &SaveRequestBody{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.EchoBindErrorTranslator(err))
	}
	if err := validator.Validate.Struct(body); err != nil {
		errs := err.(gpgvalidator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(validator.Trans))
	}

	userId, err := primitive.ObjectIDFromHex(body.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	thread := entity.Thread{
		Text:        body.Text,
		UserId:      userId,
		Likes:       body.Likes,
		RepostCount: body.RepostCount,
	}

	err = t.timelineUseCase.Save(thread)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}
