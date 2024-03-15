package email

import (
	"gopkg.in/gomail.v2"
)

type GoMailOption struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GoMail struct {
	D      *gomail.Dialer
	Option *GoMailOption
}

func NewGoMail(option *GoMailOption) *GoMail {
	d := gomail.NewDialer(option.Host, option.Port, option.Username, option.Password)
	return &GoMail{D: d, Option: option}
}

func (e *GoMail) SendSimpleMail(subject, content string, to ...string) error {
	m := gomail.NewMessage(gomail.SetCharset("utf-8"))
	m.SetHeader("From", e.Option.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	if err := e.D.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
