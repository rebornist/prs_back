package customTypes

type JWTKey struct {
	AccessPrivateKey  string `json:"access_private_key"`
	AccessPublicKey   string `json:"access_public_key"`
	RefreshPrivateKey string `json:"refresh_private_key"`
	RefreshPublicKey  string `json:"refresh_public_key"`
}

type JWTUserInfo struct {
	UID      string  `json:"uid"`
	Username string  `json:"username"`
	Grade    float64 `json:"grade"`
	IdToken  string  `json:"idToken"`
}

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessId     string `json:"access_id"`
	RefreshId    string `json:"refresh_id"`
	AtExpires    int64  `json:"access_expires"`
	RtExpires    int64  `json:"refresh_expires"`
}
