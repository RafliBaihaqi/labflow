package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuotationRepo struct {
	db *gorm.DB
}

func NewQuotationRepo(db *gorm.DB) *QuotationRepo {
	return &QuotationRepo{db: db}
}

func (r *QuotationRepo) CreateQuotation(ctx context.Context, email, companyName string, phoneNumber, message *string) (*models.Quotation, error) {
	quotation := models.Quotation{
		XID:         uuid.New(),
		Email:       email,
		CompanyName: companyName,
		PhoneNumber: phoneNumber,
		Message:     message,
	}

	if err := r.db.WithContext(ctx).Create(&quotation).Error; err != nil {
		return nil, err
	}

	return &quotation, nil
}

func (r *QuotationRepo) GetQuotations(ctx context.Context, search string, page, limit int) ([]models.Quotation, int64, error) {
	var quotations []models.Quotation
	var total int64
	query := r.db.WithContext(ctx).Model(&models.Quotation{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("email ILIKE ? OR company_name ILIKE ? OR phone_number ILIKE ?", like, like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&quotations).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}
	return quotations, total, nil
}

func (r *QuotationRepo) UpdateQuotationStatus(ctx context.Context, xid uuid.UUID, status string) (*models.Quotation, error) {
	status = strings.TrimSpace(status)
	if err := r.db.WithContext(ctx).Model(&models.Quotation{}).Where("xid = ?", xid).Update("status", status).Error; err != nil {
		return nil, fmt.Errorf("failed to update quotation: %w", wrapDBError(err))
	}
	var quotation models.Quotation
	if err := r.db.WithContext(ctx).Where("xid = ?", xid).First(&quotation).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return &quotation, nil
}
