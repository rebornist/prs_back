package file

import (
	"net/http"
	"prs/middlewares"
	"prs/models"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func OpinUpdate(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var opin models.Opin
	var resp OpinResponse

	opin.RecordId = c.FormValue("record_id")
	opin.FileId = c.FormValue("file_id")
	opin.FileNm = c.FormValue("file_nm")
	opin.NlpId = c.FormValue("nlp_id")
	opin.RecordTitle = c.FormValue("record_title")
	opin.Title = c.FormValue("title")
	opin.ReclssDivCd = c.FormValue("reclss_div_cd")
	opin.ReclssOpenGrade = c.FormValue("reclss_open_grade")
	opin.DivPart = c.FormValue("div_part")
	opin.DivInfo = c.FormValue("div_info")
	opin.Content = c.FormValue("content")
	opin.Report = c.FormValue("report")
	opin.Remarks = c.FormValue("remarks")
	opin.UserId = c.FormValue("user_id")

	docType, err := strconv.ParseInt(c.FormValue("doc_type"), 10, 64)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	permission, err := strconv.ParseInt(c.FormValue("permission"), 10, 64)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	complate, err := strconv.ParseInt(c.FormValue("complate"), 10, 64)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	opin.DocType = uint8(docType)
	opin.Permission = uint8(permission)
	opin.Complate = uint8(complate)
	opin.UpdatedAt = time.Now()

	fileId := c.Param("file_id")
	recordId := c.Param("record_id")

	if err := db.
		Table("opins").
		Order("file_id").
		Where("file_id = ? ", fileId).
		Save(&opin).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := db.Table("opins").Order("file_id").Select("file_id").Where("record_id = ?", recordId).Scan(&resp.Files).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp.Code = http.StatusOK
	resp.Message = "opin's data update successful"
	resp.Opin = opin

	middlewares.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, resp)
}
