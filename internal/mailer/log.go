package mailer

import (
	"context"
	"log"
)

// LogMailer is the dev-mode stand-in when SMTP isn't configured — logs
// instead of sending, so OTP sign-in and user invites still work locally.
type LogMailer struct{}

func (LogMailer) SendOTP(_ context.Context, email, code string) error {
	log.Printf("DEBUG: OTP code for %s: %s (SMTP not configured)", email, code)
	return nil
}

func (LogMailer) SendInviteEmail(_ context.Context, email, name string) error {
	log.Printf("DEBUG: invite email for %s <%s> (SMTP not configured)", name, email)
	return nil
}
