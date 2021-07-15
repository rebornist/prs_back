package record

import (
	"fmt"
	"net/http"
	"prs/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ConditionView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	type conditionType struct {
		ProcDepNm []string `json:"proc_dep_nm"`
		ClssNm    []string `json:"clss_nm"`
	}

	var procDepNms []string
	var clssNms []string
	var resp conditionType

	if err := db.
		Table("records").
		Order("proc_dep_nm").
		Select("proc_dep_nm").
		Group("proc_dep_nm").
		Scan(&procDepNms).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(procDepNms) == 0 {
		err := fmt.Errorf("%s", "생산기관명이 출력되지 않습니다.")
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if err := db.
		Table("records").
		Order("clss_nm").
		Select("clss_nm").
		Group("clss_nm").
		Scan(&clssNms).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if len(procDepNms) == 0 {
		err := fmt.Errorf("%s", "기록물 유형이 출력되지 않습니다.")
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp.ProcDepNm = procDepNms
	resp.ClssNm = clssNms

	middlewares.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, resp)
}
