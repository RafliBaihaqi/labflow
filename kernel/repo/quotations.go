package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"labflow/sqlc"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuotationRepo struct {
	q *sqlc.Queries
}

func NewQuotationRepo(pool *pgxpool.Pool) *QuotationRepo {
	return &QuotationRepo{q: sqlc.New(pool)}
}

func (r *QuotationRepo) CreateQuotation(ctx context.Context, email, companyName string, phoneNumber, message *string) (*models.Quotation, error) {
	row, err := r.q.CreateQuotation(ctx, sqlc.CreateQuotationParams{
		Xid:         uuid.New(),
		Email:       email,
		CompanyName: companyName,
		PhoneNumber: phoneNumber,
		Message:     message,
	})
	if err != nil {
		return nil, err
	}
	quotation := toQuotation(row)
	return &quotation, nil
}

func (r *QuotationRepo) GetQuotations(ctx context.Context, search string, page, limit int) ([]models.Quotation, int64, error) {
	pattern := likePattern(search)
	total, err := r.q.CountQuotations(ctx, pattern)
	if err != nil {
		return nil, 0, wrapDBError(err)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	rows, err := r.q.ListQuotations(ctx, sqlc.ListQuotationsParams{
		Search:     pattern,
		PageLimit:  limit,
		PageOffset: (page - 1) * limit,
	})
	if err != nil {
		return nil, 0, wrapDBError(err)
	}
	quotations := make([]models.Quotation, len(rows))
	for i, row := range rows {
		quotations[i] = toQuotation(row)
	}
	return quotations, total, nil
}

func (r *QuotationRepo) UpdateQuotationStatus(ctx context.Context, xid uuid.UUID, status string) (*models.Quotation, error) {
	status = strings.TrimSpace(status)
	row, err := r.q.UpdateQuotationStatus(ctx, sqlc.UpdateQuotationStatusParams{Xid: xid, Status: &status})
	if err != nil {
		return nil, fmt.Errorf("failed to update quotation: %w", wrapDBError(err))
	}
	quotation := toQuotation(row)
	return &quotation, nil
}
