package main

import (
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
)

type emailSender struct {
	config Config
}

type Config struct {
	userMail string
	password string

	host       string
	serverAddr string
}

func NewSender(config Config) emailSender {
	return emailSender{config: config}
}

func (es *emailSender) SendEmail() error {
	e := &email.Email{
		To:      []string{"robinsonmu232@gmail.com"},
		From:    "Robinson Muñoz Muñoz <testrobinson98@gmail.com>",
		Subject: "Awesome Subject",
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	err := e.Send(
		es.config.serverAddr,
		smtp.PlainAuth("", es.config.userMail, es.config.password, es.config.host))
	if err != nil {
		return err
	}

	return nil
}
