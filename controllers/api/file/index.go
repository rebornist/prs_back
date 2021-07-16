package file

import "github.com/labstack/echo/v4"

func OpinRouter(e *echo.Group) {
	opin := e.Group("/record")
	opin.GET("/:record_id/:file_id", DetailView)
	opin.PUT("/:record_id/:file_id", OpinUpdate)
}
