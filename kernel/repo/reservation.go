package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationRepo struct {
	db *gorm.DB
}

func NewReservationRepo(db *gorm.DB) *ReservationRepo {
	return &ReservationRepo{db: db}
}

type UpsertReservationParams struct {
	UserID              int
	SampleName          string
	SampleAmount        float64
	ReceivedNote        *string
	TestRequestNumber   *string
	SampleResidue       string
	SpecialInstructions *string
	POURL               *string
	PONumber            *string
	SampleType          string
	SamplePhoto         *string
	StartTestDate       *time.Time
	EndTestDate         *time.Time
	Note                *string
	SAPNotes            *string
	BatchNumber         *string
	SampleTakenDate     time.Time
	SampleTakenBy       string
	IsKan               *bool
	Tests               []models.TestItem
	UpdatedBy           *int
}

type UpdateStatusParams struct {
	Status      string
	UpdatedBy   *int
	SamplePhoto *string
}

func (r *ReservationRepo) CreateReservationWithTests(ctx context.Context, userID int, sampleName string, sampleAmount float64, receivedNote *string, sampleResidue string, specialInstructions *string, sampleType string, note *string, batchNumber *string, sampleTakenDate time.Time, sampleTakenBy string, testRequestNumber *string, tests []models.TestItem) (*models.Reservation, error) {
	params := UpsertReservationParams{
		UserID:              userID,
		SampleName:          sampleName,
		SampleAmount:        sampleAmount,
		ReceivedNote:        receivedNote,
		SampleResidue:       sampleResidue,
		SpecialInstructions: specialInstructions,
		SampleType:          sampleType,
		Note:                note,
		BatchNumber:         batchNumber,
		SampleTakenDate:     sampleTakenDate,
		SampleTakenBy:       sampleTakenBy,
		TestRequestNumber:   testRequestNumber,
		Tests:               tests,
	}

	reservation := models.Reservation{
		XID:                 uuid.New(),
		UserID:              params.UserID,
		ReservationDate:     time.Now(),
		SampleName:          params.SampleName,
		SampleAmount:        params.SampleAmount,
		ReceivedNote:        params.ReceivedNote,
		SampleResidue:       params.SampleResidue,
		SpecialInstructions: params.SpecialInstructions,
		SampleType:          params.SampleType,
		Note:                params.Note,
		BatchNumber:         params.BatchNumber,
		SampleTakenDate:     params.SampleTakenDate,
		SampleTakenBy:       params.SampleTakenBy,
		TestRequestNumber:   params.TestRequestNumber,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&reservation).Error; err != nil {
			return fmt.Errorf("failed to create reservation: %w", wrapDBError(err))
		}

		tests := make([]models.Test, len(params.Tests))

		for i, t := range params.Tests {
			tests[i] = models.Test{
				XID:           uuid.New(),
				ReservationID: reservation.ID,
				TestListID:    t.TestListID,
				Method:        t.Method,
				Criteria:      t.Criteria,
				TestUnit:      t.TestUnit,
			}
		}

		if len(tests) > 0 {
			if err := tx.Create(&tests).Error; err != nil {
				return fmt.Errorf("failed to create tests: %w", wrapDBError(err))
			}
		}
		reservation.Tests = tests
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

var allowedReservationSort = map[string]string{
	"test_request_number": "test_request_number",
	"sample_name":         "sample_name",
	"status":              "status",
	"reservation_date":    "reservation_date",
}

func (r *ReservationRepo) GetAllReservations(ctx context.Context, search string, status []string, page, limit int, sortBy, sortOrder string) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Reservation{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("sample_name ILIKE ? OR test_request_number ILIKE ?", like, like)
	}
	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}

	orderColumn, ok := allowedReservationSort[sortBy]
	if !ok {
		orderColumn = "created_at"
	}
	orderDir := "DESC"
	if strings.ToLower(sortOrder) == "asc" {
		orderDir = "ASC"
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Preload("Tests.TestList").
		Order(fmt.Sprintf("%s %s", orderColumn, orderDir)).
		Limit(limit).
		Offset(offset).
		Find(&reservations).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}

	return reservations, total, nil
}

func (r *ReservationRepo) GetReservationStats(ctx context.Context) (map[string]int64, int64, error) {
	type statusCount struct {
		Status *string
		Count  int64
	}
	var rows []statusCount
	if err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, 0, wrapDBError(err)
	}

	counts := make(map[string]int64)
	var total int64
	for _, row := range rows {
		key := ""
		if row.Status != nil {
			key = *row.Status
		}
		counts[key] = row.Count
		total += row.Count
	}
	return counts, total, nil
}

func (r *ReservationRepo) UpdateReservationByXid(ctx context.Context, xid uuid.UUID, poURL, poNumber, samplePhoto *string, updatedBy int) (*models.Reservation, error) {
	params := UpsertReservationParams{
		POURL:       poURL,
		PONumber:    poNumber,
		SamplePhoto: samplePhoto,
		UpdatedBy:   &updatedBy,
	}

	updates := map[string]interface{}{
		"po_url":       params.POURL,
		"po_number":    params.PONumber,
		"sample_photo": params.SamplePhoto,
		"updated_at":   time.Now(),
		"updated_by":   params.UpdatedBy,
	}

	if err := r.db.Model(&models.Reservation{}).Where("xid = ?", xid).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update reservation: %w", wrapDBError(err))
	}

	var updated models.Reservation
	if err := r.db.Where("xid = ?", xid).First(&updated).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return &updated, nil
}

func (r *ReservationRepo) GetReservationByXid(ctx context.Context, xid uuid.UUID) (*models.Reservation, error) {
	var reservation *models.Reservation
	if err := r.db.WithContext(ctx).Preload("User").Preload("Tests.TestList.Category").Where("xid = ?", xid).First(&reservation).Error; err != nil {
		return nil, wrapDBError(err)
	}
	return reservation, nil
}

func (r *ReservationRepo) UpdateReservationStatus(ctx context.Context, xid uuid.UUID, status string, updatedBy int, samplePhoto *string) (*models.Reservation, error) {
	params := UpdateStatusParams{
		Status:      status,
		UpdatedBy:   &updatedBy,
		SamplePhoto: samplePhoto,
	}

	updates := map[string]interface{}{
		"status":     params.Status,
		"updated_by": params.UpdatedBy,
	}

	if params.SamplePhoto != nil {
		updates["sample_photo"] = params.SamplePhoto
	}

	if err := r.db.WithContext(ctx).Model(&models.Reservation{}).Where("xid = ?", xid).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update reservation status: %w", wrapDBError(err))
	}

	var updated models.Reservation
	if err := r.db.WithContext(ctx).Where("xid = ?", xid).First(&updated).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return &updated, nil
}
