package mailer

import (
	"bytes"
	"html/template"
	"os"
	"strconv"

	"embed"

	"github.com/labstack/echo"
	"gopkg.in/gomail.v2"
)

//go:embed templates
var templateFS embed.FS

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

type Mailer struct {
	dialer *gomail.Dialer
	sender string
	logger echo.Logger
}

func NewMailer(logger echo.Logger) Mailer {

	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal("Invalid MAIL_PORT value")
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSender := os.Getenv("MAIL_SENDER")
	dialer := gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)
	return Mailer{
		dialer: dialer,
		sender: mailSender,
		logger: logger,
	}
}

func (mailer *Mailer) SendEmail(recipient string, templateFile string, data EmailData) error {

	absolutePath := "templates/" + templateFile
	tmpl, err := template.ParseFS(templateFS, absolutePath)
	if err != nil {
		mailer.logger.Error("Error parsing email template: ", err)
		return err
	}

	data.AppName = os.Getenv("APP_NAME")
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		mailer.logger.Error("Error executing email subject template: ", err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "body", data)
	if err != nil {
		mailer.logger.Error("Error executing email subject template:", err)
		return err
	}
	message := gomail.NewMessage()
	message.SetHeader("From", mailer.sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(message)
	if err != nil {
		mailer.logger.Error("Error sending email:", err)
		return err
	}
	return nil
}
