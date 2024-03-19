package main

import (
	"html/template"
	"log"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"

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

type SummaryMailData struct {
	TotalBalance float64

	TransactionsByMonth string

	AverageDebit  string
	AverageCredit string
}

func (es *emailSender) SendEmailNotification() error {
	// Read the HTML template file
	templatePath := "./templates/AccountSummary.html"
	htmlContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Error reading HTML template file: %v", err)
	}

	summary := SummaryMailData{}

	// Parse and execute the HTML template with data
	tmpl := template.Must(template.New("account_summary").Parse(string(htmlContent)))
	var filledTemplateContent strings.Builder
	if err := tmpl.Execute(&filledTemplateContent, summary); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	// Set the HTML content of the email with the embedded logo

	e := email.NewEmail()
	e.Subject = "Account transactions summary"
	e.From = "Robinson Muñoz Muñoz <testrobinson98@gmail.com>"
	e.To = []string{"robinsonmu232@gmail.com"}
	e.Headers = textproto.MIMEHeader{
		"Content-Type": {"text/html", "charset=utf-8"},
	}

	e.HTML = []byte(filledTemplateContent.String())

	err = e.Send(
		es.config.serverAddr,
		smtp.PlainAuth("", es.config.userMail, es.config.password, es.config.host))
	if err != nil {
		return err
	}

	return nil
}
