package common

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"prs/configs"
	"prs/customTypes"
	"prs/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateAccessJWT(username string) (*customTypes.JWTToken, *customTypes.JWTUserInfo, error) {

	jt, userInfo, err := initJWT(username)
	if err != nil {
		return nil, nil, err
	}

	// rsa 파일 위치 불러오기
	var jwtKey customTypes.JWTKey
	pathByte, err := configs.GetServiceInfo("jwt_token")
	if err != nil {
		return nil, nil, err
	}
	json.Unmarshal(pathByte, &jwtKey)

	aSignKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(jwtKey.AccessPrivateKey))
	if err != nil {
		return nil, nil, err
	}

	accessToken := jwt.New(jwt.SigningMethodRS256)
	atClaims := accessToken.Claims.(jwt.MapClaims)
	atClaims["authorized"] = true
	atClaims["user"] = userInfo
	atClaims["access_id"] = jt.AccessId
	atClaims["exp"] = jt.AtExpires
	jt.AccessToken, err = accessToken.SignedString(aSignKey)
	if err != nil {
		return nil, nil, err
	}

	return jt, userInfo, nil
}

func CreateRefreshJWT(username string) (*customTypes.JWTToken, *customTypes.JWTUserInfo, error) {

	jt, userInfo, err := initJWT(username)
	if err != nil {
		return nil, nil, err
	}

	// rsa 파일 위치 불러오기
	var jwtKey customTypes.JWTKey
	pathByte, err := configs.GetServiceInfo("jwt_token")
	if err != nil {
		return nil, nil, err
	}
	json.Unmarshal(pathByte, &jwtKey)

	rSignKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(jwtKey.RefreshPrivateKey))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodRS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["refresh_id"] = jt.RefreshId
	rtClaims["user"] = userInfo
	rtClaims["exp"] = jt.RtExpires
	jt.RefreshToken, err = refreshToken.SignedString(rSignKey)
	if err != nil {
		return nil, nil, err
	}

	return jt, userInfo, nil
}

func VerifyJWT(Type, Value string) (*jwt.Token, error) {

	// rsa 파일 위치 불러오기
	var jwtKey customTypes.JWTKey
	var pathByte, _ = configs.GetServiceInfo("jwt_token")
	json.Unmarshal(pathByte, &jwtKey)

	var verifyKey *rsa.PublicKey
	var err error
	switch Type {
	case "access":
		verifyKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(jwtKey.AccessPublicKey))
		if err != nil {
			return nil, err
		}
	case "refresh":
		verifyKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(jwtKey.RefreshPublicKey))
		if err != nil {
			return nil, err
		}
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func CheckRefreshToken(c echo.Context, db *gorm.DB) (*customTypes.JWTUserInfo, error) {
	session, err := c.Cookie("SID")
	if err != nil {
		return nil, err
	}

	uid, _ := Unsigning(session.Value)
	var sess models.Session
	if err := db.Find(&sess).Where("id=?", uid).Error; err != nil {
		return nil, err
	}

	// 세션 만료시간이 지난 경우 쿠키 및 DB 삭제
	if sess.Expires < time.Now().Unix() {

		if err := db.Where("id = ?", uid).Delete(&sess).Error; err != nil {
			return nil, err
		}

		delSession := DeleteCookie("SID", "/")
		c.SetCookie(delSession)

		message := "The session expiration period has passed."
		return nil, fmt.Errorf("%s", message)
	}

	// refresh 토큰 정보 가져오기
	data, err := GetClaimsInfo("refresh", sess.TokenValue)
	if err != nil {
		return nil, err
	}

	claimsInfo := data["user"].(map[string]interface{})
	username := claimsInfo["username"].(string)

	jt, userInfo, err := CreateAccessJWT(username)
	if err != nil {
		return nil, err
	}
	userInfo.IdToken = jt.AccessToken

	return userInfo, nil
}

func GetClaimsInfo(Type, Value string) (map[string]interface{}, error) {

	token, err := VerifyJWT(Type, Value)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	claims := token.Claims
	tmp, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tmp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func initJWT(username string) (*customTypes.JWTToken, *customTypes.JWTUserInfo, error) {
	// 유저정보 추출
	var db = configs.ConnectDb()
	userInfo := &customTypes.JWTUserInfo{}

	if err := db.Table("users").Select("uid, username, grade").Where("username = ?", username).Scan(userInfo).Error; err != nil {
		return nil, nil, err
	}
	uid := *&userInfo.UID
	uid, _ = Signing(uid)

	// 토큰 생성
	jt := &customTypes.JWTToken{}

	now := time.Now() // Go Playground 에서는 항상 시각은 2009-11-10 23:00:00 +0000 UTC 에서 시작한다.

	jt.AtExpires = now.Add(time.Minute * 15).Unix()
	jt.AccessId = CreateRandomString(24)

	jt.RtExpires = now.Add(time.Hour * 24 * 7).Unix()
	var val string
	for {
		val = CreateRandomString(24)
		resp, _ := createID(db, val)
		if resp {
			break
		}
	}
	if val != "" {
		jt.RefreshId = val
		userInfo = &customTypes.JWTUserInfo{
			UID:      uid,
			Username: *&userInfo.Username,
			Grade:    *&userInfo.Grade,
		}
	} else {
		err := errors.New("토큰 아이디 생성 실패")
		return nil, nil, err
	}

	return jt, userInfo, nil
}

func createID(db *gorm.DB, val string) (bool, error) {

	var count int64
	if err := db.Table("sessions").Where("uid = ?", val).Count(&count).Error; err != nil {
		return false, err
	}

	if count == 0 {
		return true, nil
	}

	message := "해당 아이디가 존재합니다."
	return false, fmt.Errorf("%s", message)

}
