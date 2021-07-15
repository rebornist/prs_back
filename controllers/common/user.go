package common

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"prs/customTypes"
	"prs/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// UID 찾기
func FindUserID(db *gorm.DB, str string) bool {

	// 공백 문자를 받아올 시 False 반환
	if str == "" {
		return false
	}

	// 유저 타입 선언
	var user models.User

	// UID 조회 조회 후 에러 처리
	if err := db.Where("uid = ?", str).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true
		} else {
			return false
		}
	}

	// 조회 성공 시 False 반환
	return false

}

func CheckStatus(c echo.Context, db *gorm.DB) (*customTypes.JWTUserInfo, error) {

	session, err := c.Cookie("SID")
	if err != nil {
		return nil, err
	}

	uid, err := Unsigning(session.Value)
	if err != nil {
		return nil, err
	}

	var sess models.Session
	if err := db.Where("uid = ?", uid).Find(&sess).Error; err != nil {
		return nil, err
	}

	// 세션 만료시간이 지난 경우 쿠키 및 DB 삭제
	if sess.Expires < time.Now().Unix() {

		if err := db.Delete(&sess, "uid = ?", uid).Error; err != nil {
			return nil, err
		}

		delSession := DeleteCookie("SID", "/")
		c.SetCookie(delSession)

		return nil, fmt.Errorf("%s", "세션 만료 기간이 지났습니다.")
	}

	// refresh 토큰 정보 가져오기
	data, err := GetClaimsInfo("refresh", sess.TokenValue)
	if err != nil {
		return nil, err
	}

	claimsInfo := data["user"].(map[string]interface{})

	userInfo := &customTypes.JWTUserInfo{
		UID:      claimsInfo["uid"].(string),
		Username: claimsInfo["username"].(string),
		Grade:    claimsInfo["grade"].(float64),
	}

	return userInfo, nil
}

func CreateRefreshCookie(db *gorm.DB, username string) (*http.Cookie, *customTypes.JWTUserInfo, error) {
	jt, userInfo, err := CreateRefreshJWT(username)
	if err != nil {
		return nil, nil, err
	}

	// 세션 생성
	var sess models.Session
	sess.UID = jt.RefreshId
	sess.TokenValue = jt.RefreshToken
	sess.Expires = jt.RtExpires
	sess.CreatedAt = time.Now()

	if err := db.Create(&sess).Error; err != nil {
		return nil, nil, err
	}

	signValue, _ := Signing(*&jt.RefreshId)
	session := CreateCookie("SID", signValue, "/", time.Unix(*&jt.RtExpires, 0))
	return session, userInfo, nil
}
