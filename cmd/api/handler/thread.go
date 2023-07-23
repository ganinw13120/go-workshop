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

type IThread interface {
	Save(c echo.Context) error
	Get(c echo.Context) error
}

type thread struct {
	threadUseCase usecase.IThread
	logger        log.ILogger
}

func NewThread(threadUseCase usecase.IThread, logger log.ILogger) *thread {
	return &thread{
		threadUseCase: threadUseCase,
		logger:        logger,
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
	Hashtag  string  `query:"hashtag" json:"hashtag" validate:"required"`
}

type GetResponse struct {
	Data     []entity.Thread `json:"data"`
	NextPage *string         `json:"next_page"`
}

func (t thread) Get(c echo.Context) error {
	body := &GetRequest{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.EchoBindErrorTranslator(err))
	}
	if err := validator.Validate.Struct(body); err != nil {
		errs := err.(gpgvalidator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(validator.Trans))
	}
	pageSize := DefaultPageSize
	if body.PageSize != nil {
		pageSize = *body.PageSize
	}
	result, nextPage, err := t.threadUseCase.GetTimelineFromHashtag(c.Request().Context(), body.Hashtag, body.Cursor, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	response := GetResponse{
		result,
		nextPage,
	}
	return c.JSON(http.StatusOK, response)
}

func (t thread) Save(c echo.Context) error {
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

	err = t.threadUseCase.Save(c.Request().Context(), thread)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}
