package email

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
)

// Sender provides and interface for sending email
type Sender interface {
	SendMail(subject, body, recipient string) error
}

type smtpSender struct {
	host     string
	port     string
	username string
	password string
}

// GetSMTPSender returns an EmailSender that uses SMTP
func GetSMTPSender(hostEnvVar, portEnvVar, usernameEnvVar, passwordEnvVar string) Sender {
	return &smtpSender{
		host:     os.Getenv(hostEnvVar),
		port:     os.Getenv(portEnvVar),
		username: os.Getenv(usernameEnvVar),
		password: os.Getenv(passwordEnvVar),
	}
}

func (sender *smtpSender) SendMail(subject, body, recipient string) error {
	from := mail.Address{Name: "", Address: sender.username}
	to := mail.Address{Name: "", Address: recipient}
	subj := subject

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	auth := smtp.PlainAuth("", sender.username, sender.password, sender.host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         sender.host,
	}

	c, err := smtp.Dial(fmt.Sprintf("%s:%s", sender.host, sender.port))
	if err != nil {
		return fmt.Errorf("error sending email (Dial): %w", err)
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("error sending email (Auth): %w", err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return fmt.Errorf("error sending email (Mail from): %w", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		return fmt.Errorf("error sending email (Mail to): %w", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("error sending email (Data): %w", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("error sending email (Write): %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("error sending email (Close): %w", err)
	}

	c.Quit()
	return nil
}
