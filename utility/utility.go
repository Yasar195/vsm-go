package utility

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"os"

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

type WhatsAppMessage struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

type Template struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type Language struct {
	Code string `json:"code"`
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

func sendWhatsAppMessage(message WhatsAppMessage) error {

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	url := "https://graph.facebook.com/v22.0/749411591586096/messages"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("WHATSAPP_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	// Print response details (similar to curl -i)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	fmt.Printf("\nResponse Body:\n%s\n", string(body))

	// Check for success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Message sent successfully!")
	} else {
		return fmt.Errorf("API returned error status: %s", resp.Status)
	}

	return nil
}
