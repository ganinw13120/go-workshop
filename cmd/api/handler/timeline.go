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

const DefaultPageSize = 10

type ITimeline interface {
	Save(c echo.Context) error
	Get(c echo.Context) error
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
	Text         string  `json:"text" validate:"required"`
	UserId       string  `json:"user_id" validate:"required"`
	Likes        int     `json:"likes" validate:"required"`
	ParentThread *string `json:"parent_thread"`
}

type GetRequest struct {
	Cursor   *string `query:"cursor"`
	PageSize *int    `query:"page_size"`
	Hashtag  string  `query:"hashtag"`
}

type GetResponse struct {
	Data     []entity.Thread `json:"data"`
	NextPage *string         `json:"next_page"`
}

func (t timeline) Get(c echo.Context) error {
	body := &GetRequest{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.EchoBindErrorTranslator(err))
	}

	pageSize := DefaultPageSize
	if body.PageSize != nil {
		pageSize = *body.PageSize
	}
	result, nextPage, err := t.timelineUseCase.GetTimelineFromHashtag(c.Request().Context(), body.Hashtag, body.Cursor, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	response := GetResponse{
		result,
		nextPage,
	}
	return c.JSON(http.StatusOK, response)
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

	var parentThread *primitive.ObjectID
	if body.ParentThread != nil {
		parent, err := primitive.ObjectIDFromHex(*body.ParentThread)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		parentThread = &parent
	}

	thread := entity.Thread{
		Text:         body.Text,
		UserId:       userId,
		Likes:        body.Likes,
		ParentThread: parentThread,
	}

	err = t.timelineUseCase.Save(c.Request().Context(), thread)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}
