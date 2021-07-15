package auth

import "github.com/labstack/echo/v4"

func AuthRouter(e *echo.Group) {
	auth := e.Group("/auth")
	auth.GET("", LoginView)
	auth.GET("/status", Status)
	auth.POST("/signin", PostLogin)
	auth.POST("/signup", CreateUser)
}
