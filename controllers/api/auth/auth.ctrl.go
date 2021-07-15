package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"prs/controllers/common"
	"prs/customTypes"
	"prs/middlewares"
	"prs/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Status(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	var resp customTypes.Response
	var authData customTypes.Auth

	userInfo, err := common.CheckStatus(c, db)
	if err != nil {
		respData, err := json.Marshal(authData)
		if err != nil {
			middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		resp.Code = http.StatusUnauthorized
		resp.Message = "Non User"
		resp.Data = string(respData)
		middlewares.CreateLogger(db, logger, resp.Code, err)
		return c.JSON(resp.Code, resp)
	}

	respData, err := json.Marshal(userInfo)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp.Code = http.StatusOK
	resp.Message = "signin sucessful"
	resp.Data = string(respData)

	return c.JSON(http.StatusOK, resp)
}

func LoginView(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	cookie, err := c.Cookie("_csrf")
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	middlewares.CreateLogger(db, logger, http.StatusOK, nil)
	return c.HTML(http.StatusOK, fmt.Sprintf("<input type=hidden id=%s name=%s value=%s />", cookie.Name, cookie.Name, cookie.Value))
}

func PostLogin(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	// 응답 타입 및 유저 타입 선언
	var resp customTypes.Response
	// var records []models.SearchRecord
	var user models.User

	// 요청 데이터 선언
	username := c.FormValue("username")
	password := c.FormValue("password")

	// User 조회
	if err := db.Where("username = ?", username).Find(&user).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)

	}

	if user.UID == "" {
		err := fmt.Errorf("%s", "아이디 혹은 패스워드를 확인해주세요.")
		middlewares.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return c.JSON(http.StatusUnauthorized, err)
	}

	// 패스워드 복호화
	orgPwd, err := common.Unsigning(user.Password)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 패스워드 비교 후 실패 시 에러 반환
	if orgPwd != password {
		err = fmt.Errorf("%s", "아이디 혹은 패스워드를 확인해주세요.")
		middlewares.CreateLogger(db, logger, http.StatusUnauthorized, err)
		return c.JSON(http.StatusUnauthorized, err)
	}

	sess, userInfo, err := common.CreateRefreshCookie(db, username)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.SetCookie(sess)

	respData, err := json.Marshal(userInfo)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	resp.Code = http.StatusOK
	resp.Message = "signin sucessful"
	resp.Data = string(respData)

	return c.JSON(http.StatusOK, resp)
}

func CreateUser(c echo.Context) error {
	db := c.Request().Context().Value("DB").(*gorm.DB)
	logger := c.Request().Context().Value("LOG").(*logrus.Entry)

	// 응답 타입 및 유저 타입 선언
	var resp customTypes.Response
	var data models.ResponseUser
	var user models.User

	// 요청 데이터 선언
	username := c.FormValue("username")
	password := c.FormValue("password")

	// 패스워드 변환
	convPwd, err := common.Signing(password)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// UID 생성
	uid := common.CreateRandomString(16)
	for {
		result := common.FindUserID(db, uid)
		if !result {
			break
		}
		uid = common.CreateRandomString(16)
	}

	// 요청 데이터 DB 업로드
	user.Username = username
	user.Password = convPwd
	user.UID = uid

	if err := db.Create(&user).Error; err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// USER 데이터 변환
	data.Grade = user.Grade
	data.UID = user.UID
	data.Username = user.Username

	respData, err := json.Marshal(data)
	if err != nil {
		middlewares.CreateLogger(db, logger, http.StatusInternalServerError, err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// 응답 데이터 생성
	resp.Code = http.StatusOK
	resp.Message = "create user sucessful"
	resp.Data = string(respData)

	return c.JSON(http.StatusOK, resp)
}
