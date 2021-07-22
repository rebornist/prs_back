package record

import (
	"encoding/json"
	"fmt"
	"net/http"
	"prs/middlewares"
	"prs/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ListView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	type Params struct {
		RecordTitle   string `json:"record-title"`
		FolderTitle   string `json:"folder-title"`
		OrgCd         string `json:"org-cd"`
		RecordType    string `json:"record-type"`
		OrgPrivNm     string `json:"org-priv-nm"`
		WorkComplate  string `json:"work-complate"`
		SortCondition string `json:"sort-condition"`
		ViewNumber    uint   `json:"view-number"`
		Page          uint   `json:"page"`
	}

	var records []models.Record
	var resp RecordsResponse
	var params Params
	var condition = map[string]string{
		"title":        "",
		"folder_title": "",
		"proc_dep_nm":  "",
		"clss_nm":      "",
		"open_grade":   "",
		"complate":     "",
	}
	var cnt int64
	var page int64

	orderBy := "folder_id, record_id, title"

	// userInfo, _ := common.CheckStatus(c, db)
	// _, err := common.Unsigning(userInfo.UID)
	// if err != nil {
	// 	middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
	// 	return c.JSON(http.StatusInternalServerError, err)
	// }

	if c.QueryParam("record-title") != "" {
		params.RecordTitle = c.QueryParam("record-title")
		condition["title"] = c.QueryParam("record-title")
	}
	if c.QueryParam("folder-title") != "" {
		params.FolderTitle = c.QueryParam("folder-title")
		condition["folder_title"] = c.QueryParam("folder-title")
	}
	if c.QueryParam("org-cd") != "" {
		params.OrgCd = c.QueryParam("org-cd")
		condition["proc_dep_nm"] = c.QueryParam("org-cd")
	}
	if c.QueryParam("record-type") != "" {
		params.RecordType = c.QueryParam("record-type")
		condition["clss_nm"] = c.QueryParam("record-type")
	}
	if c.QueryParam("org-priv-nm") != "" {
		params.OrgPrivNm = c.QueryParam("org-priv-nm")
		condition["open_grade"] = c.QueryParam("org-priv-nm")
	}
	if c.QueryParam("work-complate") != "" {
		params.WorkComplate = c.QueryParam("work-complate")
		condition["complate"] = c.QueryParam("work-complate")
	}
	if c.QueryParam("sort-condition") != "" {
		orderBy = strings.Replace(c.QueryParam("sort-condition"), "--", " ", -1)
		params.SortCondition = c.QueryParam("sort-condition")
	}

	pageNumber := 10
	if c.QueryParam("view-number") != "" {
		gpn, err := strconv.ParseInt(c.QueryParam("view-number"), 10, 64)
		if err != nil {
			middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		pageNumber = int(gpn)
	}

	params.ViewNumber = uint(pageNumber)

	if c.QueryParam("page") != "" {
		cp, err := strconv.ParseInt(c.QueryParam("page"), 10, 64)
		if err != nil {
			middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		page = cp
	} else {
		page = 1
	}

	params.Page = uint(page)

	var offset int
	if page <= 1 {
		offset = 0
	} else {
		offset = (int(page) - 1) * pageNumber
	}

	q := "complate >= ?"
	for k, v := range condition {
		if v != "" {
			if k == "title" || k == "folder_title" {
				q = fmt.Sprintf("%s AND %s LIKE '%s%s%s'", q, k, "%", v, "%")
			} else if k == "open_grade" {
				cvi, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
					return c.JSON(http.StatusInternalServerError, err)
				}
				var g string
				for i := 0; i < int(cvi); i++ {
					if i == int(cvi)-1 {
						g = fmt.Sprintf("%s%s", g, "Y")
					} else {
						g = fmt.Sprintf("%s%s", g, "_")
					}
				}
				q = fmt.Sprintf("%s AND %s LIKE '%s%s'", q, k, g, "%")
			} else {
				q = fmt.Sprintf("%s AND %s = '%s'", q, k, v)
			}
		}
	}

	if err := db.Table("records").Where(q, 0).Count(&cnt).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if int(cnt) != 0 && int(cnt) < offset {
		err := fmt.Errorf("%s", "입력된 페이지 값이 올바르지 않습니다.")
		middlewares.CreateLogger(db, logger, http.StatusBadRequest, err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"code": http.StatusBadRequest, "message": err.Error()})
	}

	if err := db.
		Table("records").
		Order(orderBy).
		Limit(pageNumber).
		Offset(offset).
		Select("record_id, proc_dep_nm, clss_nm, folder_title, title, open_div_cd, open_grade, complate").
		Where(q, 0).
		Scan(&records).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp.Message = "search records successful"
	resp.Records = records
	if len(resp.Records) == 0 {
		resp.Page = 0
		resp.TotalItems = 0
	} else {
		resp.Page = int(page)
		resp.TotalItems = int(cnt)
	}
	bytxt, _ := json.Marshal(params)
	json.Unmarshal(bytxt, &resp.QueryParam)
	resp.ViewNumber = int(pageNumber)

	middlewares.CreateLogger(db, logger, http.StatusOK, nil)
	return c.JSON(http.StatusOK, resp)
}
