package queue

import (
	"fmt"
	"labflow/mail"
	"labflow/worker"

	"github.com/hibiken/asynq"
)

// Enqueuer wraps an asynq.Client and exposes domain-level enqueue methods.
// It is the only type that touches asynq.Client; services depend on the
// narrow interface defined in their own package.
type Enqueuer struct {
	client *asynq.Client
}

// NewEnqueuer constructs an Enqueuer from a full Redis URL (e.g. redis://host:6379/0).
func NewEnqueuer(redisURL string) *Enqueuer {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		panic("queue: invalid REDIS_URL for asynq: " + err.Error())
	}
	client := asynq.NewClient(opt)
	return &Enqueuer{client: client}
}

// Close releases the underlying asynq client connection.
func (e *Enqueuer) Close() error {
	return e.client.Close()
}

// SendQuotationRequestEmail satisfies the service.quotationMailer interface.
// It enqueues the email task instead of sending it inline.
func (e *Enqueuer) SendQuotationRequestEmail(to string, data mail.QuotationRequestEmailData) error {
	task, err := worker.NewSendQuotationEmailTask(to, data)
	if err != nil {
		return fmt.Errorf("queue: create quotation email task: %w", err)
	}
	if _, err := e.client.Enqueue(task); err != nil {
		return fmt.Errorf("queue: enqueue quotation email task: %w", err)
	}
	return nil
}

func (e *Enqueuer) SendVerificationEmail(to string, otp string) error {
	task, err := worker.NewSendOTPEmailTask(to, otp)
	if err != nil {
		return fmt.Errorf("queue: create otp email task: %w", err)
	}
	if _, err := e.client.Enqueue(task); err != nil {
		return fmt.Errorf("queue: enqueue otp email task: %w", err)
	}
	return nil
}
