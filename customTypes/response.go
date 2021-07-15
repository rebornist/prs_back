package customTypes

// 공통 Responst type
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
