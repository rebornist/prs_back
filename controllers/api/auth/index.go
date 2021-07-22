package auth

import "github.com/labstack/echo/v4"

func AuthRouter(e *echo.Group) {
	auth := e.Group("/auth")
	auth.GET("/status", Status)
	auth.GET("/signin", LoginView)
	auth.POST("/signin", PostLogin)
	auth.POST("/signup", CreateUser)
}
