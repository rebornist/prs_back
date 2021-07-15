package main

import (
	"net/http"
	"prs/configs"
	"prs/middlewares"

	"prs/controllers/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := configs.ConnectDb()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// 각 request마다 고유의 ID를 부여
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middlewares.DbContext(db))
	e.Use(middlewares.LogrusLogger())

	api.AppRouter(e)

	e.Logger.Fatal(e.Start(":1323"))
}
