package emailcontroller

import (
	"fmt"
	"log"

	"github.com/upload-media-email/config"
	"gopkg.in/gomail.v2"
)

func sendEmailHelper(to, subject, body string) error {
	// Set up the SMTP client
	smtpHost:= config.GetConfig().SMTP_HOST
	smtpPassword:= config.GetConfig().SMTP_PASSWORD
	smtpPort:= config.GetConfig().SMTP_PORT
	smtpUsername:= config.GetConfig().SMTP_USER

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", "arko@test.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func SendEmail(email string, pictureURL string) {
	subject := "Confirmation of Image Upload"
	body := fmt.Sprintf(`
		<html>
			<head></head>
			<body>
				<h1>Hi</h1>
				<p>Find the picture you uploaded in the attachment!</p><br>
				<img src="%s" alt="Broken Image">
			</body>
		</html>
	`, pictureURL)

	err := sendEmailHelper(email, subject, body)
	if err != nil {
		log.Println("error while sending mail :", err.Error())
	}

	log.Println("mail sent successfully")
}