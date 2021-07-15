package file

import (
	"fmt"
	"net/http"
	"prs/controllers/common"
	"prs/middlewares"
	"prs/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func DetailView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var opin models.Opin
	var resp OpinResponse
	// var files []string
	fileId := c.Param("file_id")
	recordId := c.Param("record_id")

	if _, err := common.CheckStatus(c, db); err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	csrf, err := c.Cookie("_csrf")
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if fileId == "initial" {
		if err := db.
			Table("opins").
			Order("file_id").
			First(&opin).Error; err != nil {
			middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		if err := db.
			Table("opins").
			Order("file_id").
			Where("file_id = ? ", fileId).
			Scan(&opin).Error; err != nil {
			middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	if err := db.Table("opins").Order("file_id").Select("file_id").Where("record_id = ?", recordId).Scan(&resp.Files).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp.Code = http.StatusOK
	resp.Message = "opin's data search successful"
	resp.Opin = opin
	resp.Token = fmt.Sprintf("<input type=hidden id=%s name=%s value=%s />", csrf.Name, csrf.Name, csrf.Value)

	middlewares.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, resp)
}
