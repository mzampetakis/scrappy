package informer

import (
	"fmt"
	"net/smtp"
)

type InformUser interface {
	Inform(string, string) error
}

type MailConfig struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (mailConfig *MailConfig) Inform(title string, message string) error {
	emailHost := "smtp.gmail.com"
	emailTo := mailConfig.Email
	emailPassword := mailConfig.Password
	emailPort := 587

	emailAuth := smtp.PlainAuth("", emailTo, emailPassword, emailHost)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + title + "!\n"
	msg := []byte(subject + mime + "\n" + message)
	addr := fmt.Sprintf("%s:%d", emailHost, emailPort)
	err := smtp.SendMail(addr, emailAuth, emailTo, []string{emailTo}, msg)
	return err
}
