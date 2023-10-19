package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/enzhas/feedback_back/initializers"
	"github.com/enzhas/feedback_back/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// 👇 Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *models.User, data *EmailData, emailTemp string) error {
	config, err := initializers.LoadConfig()

	if err != nil {
		log.Fatal("could not load config", err)
		return err
	}

	// Sender data.
	from := config.EmailFrom
	password := config.SMTPPass
	to := user.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort
	var body bytes.Buffer

	emailTemplate, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	if err := emailTemplate.ExecuteTemplate(&body, emailTemp, &data); err != nil {
		log.Fatal("Error loading email template")
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		//log.Fatal("Could not send email")
		return err
	}

	return nil
}
