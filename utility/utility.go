package utility

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"golang.org/x/crypto/bcrypt"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

type EmailRequest struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

type Response[T any] struct {
	Success    bool   `json:"success"`
	Data       *T     `json:"data"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"statusCode"`
}

type EmailService struct {
	config EmailConfig
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (es *EmailService) SendEmail(to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", es.config.SMTPHost, es.config.SMTPPort)

	if es.config.SMTPPort == 465 {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         es.config.SMTPHost,
		}

		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("failed to connect TLS: %w", err)
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, es.config.SMTPHost)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer c.Quit()

		auth := smtp.PlainAuth("", es.config.SMTPUsername, es.config.SMTPPassword, es.config.SMTPHost)
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}

		if err = c.Mail(es.config.FromEmail); err != nil {
			return err
		}
		if err = c.Rcpt(to); err != nil {
			return err
		}

		w, err := c.Data()
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
			es.config.FromEmail, to, subject, body)

		if _, err = w.Write([]byte(msg)); err != nil {
			return err
		}
		w.Close()

		return nil
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	c, err := smtp.NewClient(conn, es.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer c.Quit()

	tlsconfig := &tls.Config{ServerName: es.config.SMTPHost}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(tlsconfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	}

	auth := smtp.PlainAuth("", es.config.SMTPUsername, es.config.SMTPPassword, es.config.SMTPHost)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	if err = c.Mail(es.config.FromEmail); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		es.config.FromEmail, to, subject, body)

	if _, err = w.Write([]byte(msg)); err != nil {
		return err
	}
	w.Close()

	return nil
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
