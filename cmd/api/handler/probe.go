package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wisesight/go-api-template/pkg/adapter"
	"github.com/wisesight/go-api-template/pkg/log"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type IProbe interface {
	DBReadyCheck(c echo.Context) error
}

type probe struct {
	mongoDBAdapter adapter.IMongoDBAdapter
	logger         log.ILogger
}

func NewProbe(mongoDBAdapter adapter.IMongoDBAdapter, logger log.ILogger) IProbe {
	return &probe{
		mongoDBAdapter,
		logger,
	}
}

func (h probe) DBReadyCheck(c echo.Context) error {
	err := h.mongoDBAdapter.Ping(c.Request().Context(), readpref.Primary())
	if err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	return c.NoContent(http.StatusOK)
}
