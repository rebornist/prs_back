package api

import (
	"prs/controllers/api/auth"
	"prs/controllers/api/file"
	"prs/controllers/api/record"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AppRouter(e *echo.Echo) {
	api := e.Group("/api")
	api.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:  "header:X-XSRF-TOKEN",
		CookieSecure: true,
		CookiePath:   "/api",
	}))

	auth.AuthRouter(api)
	record.RecordRouter(api)
	file.OpinRouter(api)
}
