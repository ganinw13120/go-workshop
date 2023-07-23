package handler

import (
	gpgvalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/helper"
	"github.com/wisesight/go-api-template/pkg/log"
	"github.com/wisesight/go-api-template/pkg/usecase"
	"github.com/wisesight/go-api-template/pkg/validator"
	"net/http"
)

type IAccount interface {
	Save(c echo.Context) error
	Get(c echo.Context) error
}

type account struct {
	accountUseCase usecase.IAccount
	logger         log.ILogger
}

func NewAccount(accountUseCase usecase.IAccount, logger log.ILogger) *account {
	return &account{
		accountUseCase: accountUseCase,
		logger:         logger,
	}
}

type SaveAccountRequestBody struct {
	DisplayName     string `json:"display_name" validate:"required"`
	Username        string `json:"username" validate:"required"`
	ProfileImageUrl string `json:"profile_image_url" validate:"required"`
	Description     string `json:"description" validate:"required"`
	Follower        int    `json:"follower" validate:"required"`
	Following       int    `json:"following" validate:"required"`
}
type GetAccountRequest struct {
	Id string `param:"id" validate:"required"`
}

func (t account) Save(c echo.Context) error {
	body := &SaveAccountRequestBody{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.EchoBindErrorTranslator(err))
	}
	if err := validator.Validate.Struct(body); err != nil {
		errs := err.(gpgvalidator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(validator.Trans))
	}

	account := entity.Account{
		DisplayName:     body.DisplayName,
		Username:        body.Username,
		ProfileImageUrl: body.ProfileImageUrl,
		Description:     body.Description,
		Follower:        body.Follower,
		Following:       body.Following,
	}

	err := t.accountUseCase.Save(c.Request().Context(), account)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (t account) Get(c echo.Context) error {
	body := &GetAccountRequest{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.EchoBindErrorTranslator(err))
	}
	if err := validator.Validate.Struct(body); err != nil {
		errs := err.(gpgvalidator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(validator.Trans))
	}

	account, err := t.accountUseCase.Get(c.Request().Context(), body.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}
