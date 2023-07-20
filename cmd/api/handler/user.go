package handler

import (
	"fmt"
	"net/http"
	"time"

	gpgvalidator "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/constant"
	"github.com/wisesight/go-api-template/pkg/entity"
	"github.com/wisesight/go-api-template/pkg/helper"
	"github.com/wisesight/go-api-template/pkg/log"
	"github.com/wisesight/go-api-template/pkg/usecase"
	"github.com/wisesight/go-api-template/pkg/validator"
)

type IUser interface {
	GetAll(c echo.Context) error
	GetUser(c echo.Context) error
	Create(c echo.Context) error
}

type user struct {
	userUseCase usecase.IUser
	logger      log.ILogger
}

func NewUser(logger log.ILogger) IUser {
	newUserValidation()
	return &user{
		logger: logger,
	}
}

// GetAll godoc
// @id           get-all-users
// @summary      Show all users
// @description  Show all active users in the system
// @tags         users
// @accept       json
// @produce      json
// @success      200  {array}   CreateResponseBody
// @failure      400  {object}  echo.HTTPError
// @failure      404  {object}  echo.HTTPError
// @failure      500  {object}  echo.HTTPError
// @router       /users [get]
func (h user) GetAll(c echo.Context) error {
	users, err := h.userUseCase.GetAll()
	ctx := c.Request().Context()

	h.logger.Info(ctx, "get all users", log.Any("users", users))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (h user) GetUser(c echo.Context) error {
	user := c.Get(constant.JWT_CONTEXT_KEY).(entity.UserSession)

	fmt.Println(user)

	return c.JSON(http.StatusOK, user)
}

type CreateRequestBody struct {
	Name      string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"is_valid_password,required"`
	BirthDate time.Time `json:"birth_date" validate:"required"`
}

type CreateResponseBody struct {
	Name      string    `json:"name" example:"John Doe"`
	Username  string    `json:"username" example:"johndoe"`
	BirthDate time.Time `json:"birth_date" example:"2006-01-02"`
}

// Create godoc
// @id           create-user
// @summary      Create a user
// @description  Create a user in the system
// @tags         users
// @accept       json
// @produce      json
// @param  data  body  entity.User  true  "User data"
// @success      200  {object}  CreateResponseBody  "Return user data"
// @failure      400  {object}  echo.HTTPError
// @failure      404  {object}  echo.HTTPError
// @failure      500  {object}  echo.HTTPError
// @router       /users [post]
func (h user) Create(c echo.Context) error {
	body := &CreateRequestBody{}
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, helper.EchoBindErrorTranslator(err))
	}
	if err := validator.Validate.Struct(body); err != nil {
		errs := err.(gpgvalidator.ValidationErrors)
		return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(validator.Trans))
	}
	return c.JSON(http.StatusCreated, &CreateResponseBody{
		Name:      body.Name,
		Username:  body.Password,
		BirthDate: body.BirthDate,
	})
}
