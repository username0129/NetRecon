package response

type CaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	CaptchaImg    string `json:"captchaImg"`
	CaptchaLength int    `json:"captchaLength"`
	OpenCaptcha   bool   `json:"openCaptcha"`
}
