package middlewares

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db := c.Request().Context().Value("DB").(*gorm.DB)
			logger := c.Request().Context().Value("LOG").(*logrus.Entry)

			message := errors.New("에러가 발생했습니다. 관리자에게 문의하세요")

			CreateLogger(db, logger, http.StatusInternalServerError, message)
			return echo.NewHTTPError(http.StatusInternalServerError, message)
		}
	}
}
