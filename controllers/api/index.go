package api

import (
	"prs/controllers/api/auth"
	"prs/controllers/api/file"
	"prs/controllers/api/record"

	"github.com/labstack/echo/v4"
)

func AppRouter(e *echo.Echo) {
	api := e.Group("/api")

	auth.AuthRouter(api)
	record.RecordRouter(api)
	file.OpinRouter(api)
}
