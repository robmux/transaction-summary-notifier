package repositories

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	config Config
}

type Config struct {
	UserMail string
	Password string

	Host       string
	Port       int
	ServerAddr string
}

func NewSender(config Config) EmailSender {
	return EmailSender{config: config}
}

type SummaryMailData struct {
	TotalBalance float64

	TransactionsByMonth string

	AverageDebit  string
	AverageCredit string
}

func (es *EmailSender) SendEmailNotification() error {
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

	message := gomail.NewMessage()
	message.SetHeader("From", "mail_notification@stori.com")
	message.SetHeader("To", "robinsonmu232@gmail.com")
	message.SetHeader("Subject", "Account Transactions Summary")

	message.SetBody("text/html", filledTemplateContent.String())
	message.Embed("templates/resources/stori_logo.png")
	message.Attach("input/user_1_transactions.csv")

	// Connect to the SMTP server with TLS encryption
	d := gomail.NewDialer(es.config.Host, es.config.Port, es.config.UserMail, es.config.Password)
	if err := d.DialAndSend(message); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email sent successfully")

	return nil
}
