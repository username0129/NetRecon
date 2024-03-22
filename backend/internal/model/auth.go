package model

type LoginRequest struct {
	Username  string `json:"username"`   // 用户名
	Password  string `json:"password"`   // 登陆密码
	Answer    string `json:"answer"`     // 验证码
	CaptchaId string `json:"captcha_id"` // 验证码 ID
}

type LoginUser struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 登陆密码
}
