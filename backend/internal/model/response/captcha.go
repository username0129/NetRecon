package response

type CaptchaResponse struct {
	CaptchaId   string `json:"captchaId"`
	CaptchaImg  string `json:"captchaImg"`
	OpenCaptcha bool   `json:"openCaptcha"`
}
