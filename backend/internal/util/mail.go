package util

import (
	"errors"
	"net/smtp"
	"strings"
)

func SendMail(server string, port string, from string, password string, to []string, subject string, body string) error {
	// 确保收件人列表不为空
	if len(to) == 0 {
		return errors.New("收件人列表为空")
	}

	toHeader := strings.Join(to, ", ")

	// 邮件消息
	message := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"From: " + from + "\r\n" +
		"To: " + toHeader + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	// SMTP 服务认证
	auth := smtp.PlainAuth("", from, password, server)

	// 发送邮件
	err := smtp.SendMail(server+":"+port, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}
