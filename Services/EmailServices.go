package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendForgetPasswordEmail(to string, link string) error {
	// SMTP server configuration
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Sender data
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	// Email content
	subject := "Subject: Forget password link\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\r\n"
	body := "Click here to reset your password: " + link + "\r\n"

	// Combine email parts
	message := []byte(subject + mime + "\r\n" + body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}

func SendResetPasswordEmail(to string, code string) {
	// TODO: Send email
}
