package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"labflow/mail"
	"log"

	"github.com/hibiken/asynq"
)

const TypeSendQuotationEmail = "email:quotation_request"
const TypeSendOTPEmail = "email:otp_verification"

type QuotationEmailPayload struct {
	To   string
	Data mail.QuotationRequestEmailData
}

type OTPEmailPayload struct {
	To  string
	OTP string
}

func NewSendOTPEmailTask(to string, otp string) (*asynq.Task, error) {
	payload, err := json.Marshal(OTPEmailPayload{To: to, OTP: otp})
	if err != nil {
		return nil, fmt.Errorf("worker: marshal otp email payload: %w", err)
	}
	return asynq.NewTask(
		TypeSendOTPEmail,
		payload,
		asynq.MaxRetry(5),
		asynq.Queue("email"),
	), nil
}

func HandleSendOTPEmailTask(mailer *mail.Sender) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p OTPEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("worker: unmarshal quotation email payload: %w", asynq.SkipRetry)
		}

		if err := mailer.SendVerificationEmail(p.To, p.OTP); err != nil {
			log.Printf("[worker] failed to send verification email to %s: %v", p.To, err)
			return err
		}

		log.Printf("[worker] verification email sent to %s", p.To)
		return nil
	}
}

func NewSendQuotationEmailTask(to string, data mail.QuotationRequestEmailData) (*asynq.Task, error) {
	payload, err := json.Marshal(QuotationEmailPayload{To: to, Data: data})
	if err != nil {
		return nil, fmt.Errorf("worker: marshal quotation email payload: %w", err)
	}
	return asynq.NewTask(
		TypeSendQuotationEmail,
		payload,
		asynq.MaxRetry(5),
		asynq.Queue("email"),
	), nil
}

func HandleSendQuotationEmailTask(mailer *mail.Sender) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p QuotationEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("worker: unmarshal quotation email payload: %w", asynq.SkipRetry)
		}

		if err := mailer.SendQuotationRequestEmail(p.To, p.Data); err != nil {
			log.Printf("[worker] failed to send quotation email to %s: %v", p.To, err)
			return err
		}

		log.Printf("[worker] quotation email sent to %s (quotation %s)", p.To, p.Data.QuotationID)
		return nil
	}
}
