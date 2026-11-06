package service

import (
	"context"
	"errors"
	"labflow/mail"
	"labflow/models"
	"os"
	"strings"

	"github.com/google/uuid"
)

type CreateQuotationParams struct {
	Email       string
	CompanyName string
	PhoneNumber *string
	Message     *string
}

type quotationEnqueuer interface {
	SendQuotationRequestEmail(to string, data mail.QuotationRequestEmailData) error
}

type quotationRepo interface {
	CreateQuotation(ctx context.Context, email, companyName string, phoneNumber, message *string) (*models.Quotation, error)
	GetQuotations(ctx context.Context, search string, page, limit int) ([]models.Quotation, int64, error)
	UpdateQuotationStatus(ctx context.Context, xid uuid.UUID, status string) (*models.Quotation, error)
}

type QuotationService struct {
	repo   quotationRepo
	mailer quotationEnqueuer
}

// quotationInbox is where new quotation requests are emailed. Set QUOTATION_INBOX.
func quotationInbox() string {
	if v := os.Getenv("QUOTATION_INBOX"); v != "" {
		return v
	}
	return "sales@labflow.local"
}

func NewQuotationService(repo quotationRepo, mailer quotationEnqueuer) *QuotationService {
	return &QuotationService{repo: repo, mailer: mailer}
}

func (qs *QuotationService) CreateQuotation(ctx context.Context, req CreateQuotationParams) (*models.Quotation, error) {
	quotation, err := qs.repo.CreateQuotation(ctx, req.Email, req.CompanyName, req.PhoneNumber, req.Message)
	if err != nil {
		return nil, err
	}

	phoneNumber := ""
	if quotation.PhoneNumber != nil {
		phoneNumber = *quotation.PhoneNumber
	}
	message := ""
	if quotation.Message != nil {
		message = *quotation.Message
	}

	emailData := mail.QuotationRequestEmailData{
		QuotationID:   quotation.XID.String(),
		CustomerEmail: quotation.Email,
		CompanyName:   quotation.CompanyName,
		PhoneNumber:   phoneNumber,
		Message:       message,
	}

	if err := qs.mailer.SendQuotationRequestEmail(quotationInbox(), emailData); err != nil {
		return nil, err
	}

	return quotation, nil
}

var ErrInvalidQuotationStatus = errors.New("invalid quotation status")

type QuotationListResult struct {
	Data       []models.Quotation
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

func (qs *QuotationService) ListQuotations(ctx context.Context, search string, page, limit int) (QuotationListResult, error) {
	page, limit = normalPagination(page, limit)
	quotations, total, err := qs.repo.GetQuotations(ctx, search, page, limit)
	if err != nil {
		return QuotationListResult{}, err
	}
	totalPages := calculateTotalPage(total, limit)
	return QuotationListResult{Data: quotations, Total: total, Page: page, Limit: limit, TotalPages: totalPages}, nil
}

func (qs *QuotationService) UpdateQuotationStatus(ctx context.Context, xid uuid.UUID, status string) (*models.Quotation, error) {
	status = strings.TrimSpace(status)
	switch status {
	case "pending", "contacted", "closed":
	default:
		return nil, ErrInvalidQuotationStatus
	}
	return qs.repo.UpdateQuotationStatus(ctx, xid, status)
}
