package mailer

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

// SMTPMailer sends real email via net/smtp (stdlib — no library needed; it
// handles STARTTLS automatically when the server advertises it, e.g. Gmail:587).
type SMTPMailer struct {
	cfg SMTPConfig
}

func NewSMTPMailer(cfg SMTPConfig) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) SendOTP(_ context.Context, email, code string) error {
	subject := "Your sign-in code"
	body := fmt.Sprintf("Your one-time code is: %s\nIt expires in 15 minutes.", code)
	return m.send(email, subject, body)
}

func (m *SMTPMailer) SendInviteEmail(_ context.Context, email, name string) error {
	subject := "You've been invited"
	body := fmt.Sprintf("Hi %s,\n\nYou've been invited to the Horse Stable Management System.\nSign in with your email (%s) to receive a one-time code.", name, email)
	return m.send(email, subject, body)
}

func (m *SMTPMailer) send(to, subject, body string) error {
	addr := m.cfg.Host + ":" + m.cfg.Port
	auth := smtp.PlainAuth("", m.cfg.User, m.cfg.Password, m.cfg.Host)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", m.cfg.From, to, subject, body))

	if err := smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg); err != nil {
		return fmt.Errorf("send mail to %s: %w", to, err)
	}
	return nil
}

// Configured reports whether enough SMTP config is present to actually send.
func (c SMTPConfig) Configured() bool {
	return c.Host != "" && c.Port != "" && c.User != "" && c.Password != "" && c.From != ""
}
