package common

import (
	"net/http"
	"time"
)

func CreateCookie(name, token, path string, setTime time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = setTime
	// cookie.HttpOnly = true
	// cookie.Secure = true
	cookie.Path = path
	return cookie
}

func DeleteCookie(name, path string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	cookie.Path = path
	return cookie
}
