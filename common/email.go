package common

import (
	"fmt"
	"net/smtp"
	"os"

	"crypto/tls"
	"net/mail"
)

func sendHTMLMail(subject, body, recipient, hostEnvVar, portEnvVar, usernameEnvVar, passwordEnvVar string) error {

	smtpHost := os.Getenv(hostEnvVar)
	smtpPort := os.Getenv(portEnvVar)
	smtpUsername := os.Getenv(usernameEnvVar)
	smtpPassword := os.Getenv(passwordEnvVar)

	from := mail.Address{Name: "", Address: smtpUsername}
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
	servername := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	c, err := smtp.Dial(servername)
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
