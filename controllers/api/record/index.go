package record

import "github.com/labstack/echo/v4"

func RecordRouter(e *echo.Group) {
	record := e.Group("/records")
	record.GET("", ListView)
	record.GET("/conditions", ConditionView)
}
