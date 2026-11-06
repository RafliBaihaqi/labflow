package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

//go:embed templates/*.html
var templateFS embed.FS

type Sender struct {
	host string
	port string
	from string
	pass string
}

type QuotationRequestEmailData struct {
	QuotationID   string
	CustomerEmail string
	CompanyName   string
	PhoneNumber   string
	Message       string
}

func NewSender() *Sender {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_USER")
	if from == "" {
		from = "noreply@labflow.local"
	}

	if host == "" || port == "" {
		panic("SMTP_HOST, SMTP_PORT, SMTP_USER, and SMTP_PASS must be set")
	}

	return &Sender{
		host: host,
		port: port,
		from: from,
		pass: os.Getenv("SMTP_PASS"),
	}
}

func renderVerificationEmail(otp string) ([]byte, error) {
	t, err := template.ParseFS(templateFS, "templates/otp_verification.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse email template: %w", err)
	}
	var body bytes.Buffer
	data := struct{ OTP string }{OTP: otp}
	if err := t.Execute(&body, data); err != nil {
		return nil, fmt.Errorf("failed to execute email template: %w", err)
	}
	return body.Bytes(), nil
}

func renderQuotationRequestEmail(data QuotationRequestEmailData) ([]byte, error) {
	t, err := template.ParseFS(templateFS, "templates/quotation_request.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse email template: %w", err)
	}
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return nil, fmt.Errorf("failed to execute email template: %w", err)
	}
	return body.Bytes(), nil
}

func (s *Sender) SendVerificationEmail(to, otp string) error {
	body, err := renderVerificationEmail(otp)
	if err != nil {
		return err
	}

	fromHeader := fmt.Sprintf("From: %s\r\n", s.from)
	subject := "Subject: LabFlow Email Verification\r\n"
	toHeader := fmt.Sprintf("To: %s\r\n", to)
	mimeHeaders := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(fromHeader + subject + toHeader + mimeHeaders + string(body))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	var auth smtp.Auth
	if s.from != "" && s.pass != "" {
		auth = smtp.PlainAuth("", s.from, s.pass, s.host)
	}

	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}

func (s *Sender) SendQuotationRequestEmail(to string, data QuotationRequestEmailData) error {
	body, err := renderQuotationRequestEmail(data)
	if err != nil {
		return err
	}

	fromHeader := fmt.Sprintf("From: %s\r\n", s.from)
	subject := "Subject: New Quotation Request\r\n"
	toHeader := fmt.Sprintf("To: %s\r\n", to)
	mimeHeaders := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

	msg := []byte(fromHeader + subject + toHeader + mimeHeaders + string(body))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	var auth smtp.Auth
	if s.from != "" && s.pass != "" {
		auth = smtp.PlainAuth("", s.from, s.pass, s.host)
	}

	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
