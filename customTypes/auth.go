package customTypes

type Auth struct {
	UID        string `json:"uid"`
	Username   string `json:"username"`
	Grade      uint8  `json:"grade"`
	IsLoggedIn uint8  `json:"is_logged"`
}
