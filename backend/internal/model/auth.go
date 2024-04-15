package model

type LoginRequest struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 登陆密码
	Answer    string `json:"answer"`    // 验证码
	CaptchaId string `json:"captchaId"` // 验证码 ID
}

type LoginResponse struct {
	User  User   `json:"user"`  // 用户信息
	Token string `json:"token"` // 登陆密码
}

type LoginUser struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 登陆密码
}
