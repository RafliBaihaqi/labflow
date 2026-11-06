package service

import (
	"context"
	"errors"
	"fmt"
	"labflow/lifecycle"
	"labflow/models"
	"math"
	"time"

	"github.com/google/uuid"
)

var sampleTypeCode = map[string]string{
	"Solid":  "001",
	"Liquid": "002",
	"Powder": "003",
}

type TestItemParams struct {
	TestListID int
	Method     string
	Criteria   string
	TestUnit   string
}

type CreateReservationParams struct {
	SampleName          string
	SampleAmount        float64
	ReceivedNote        *string
	SampleResidue       string
	SpecialInstructions *string
	SampleType          string
	Note                *string
	BatchNumber         *string
	SampleTakenDate     time.Time
	SampleTakenBy       string
	Tests               []TestItemParams
}

type UpdateReservationParams struct {
	POURL       *string
	PONumber    *string
	SamplePhoto *string
}

type UpdateStatusParams struct {
	Status      string
	SamplePhoto *string
}

type ListParams struct {
	Search    string
	Status    []string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type TestResult struct {
	ID            int
	XID           uuid.UUID
	ReservationID int
	TestListId    int
	Method        string
	Criteria      string
	Status        *string
	TestUnit      string
	TestResult    *string
	TestSummary   *string
	ApprovedDate  *time.Time
	Notes         *string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	TestList      string
	CategoryName  string
}

type ReservationUser struct {
	ID          int
	XID         uuid.UUID
	Email       string
	PhoneNumber string
	Name        string
	Company     string
	Address     string
	Avatar      *string
}

type ReservationResult struct {
	ID                  int
	Xid                 uuid.UUID
	UserID              int
	ReservationDate     time.Time
	Status              *string
	SampleName          string
	SampleAmount        float64
	ReceivedNote        *string
	TestRequestNumber   *string
	SampleResidue       string
	SpecialInstructions *string
	POURL               *string
	PONumber            *string
	SampleType          string
	FeedbackRequired    *bool
	MRApproved          *bool
	HeadApproved        *bool
	MRApprovalDate      *time.Time
	HeadApprovalDate    *time.Time
	SamplePhoto         *string
	Note                *string
	StartTestDate       *time.Time
	EndTestDate         *time.Time
	SAPNotes            *string
	BatchNumber         *string
	LeadTime            *int
	SampleTakenDate     time.Time
	SampleTakenBy       string
	IsKan               *bool
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
	UpdatedBy           *int
	Message             string
	Tests               []TestResult
	User                ReservationUser
}

type CreateReservationResult struct {
	ID                  int
	Xid                 uuid.UUID
	SampleName          string
	SampleAmount        float64
	SpecialInstructions *string
	SampleType          string
	BatchNumber         *string
	SampleTakenDate     time.Time
	SampleTakenBy       string
	Message             string
	TestRequestNumber   *string
	Tests               []TestItemParams
}

type PaginatedReservationResult struct {
	Data       []ReservationResult
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

type ReservationStatsResult struct {
	Total          int64
	Pending        int64
	SampleReceived int64
	TestOnProgress int64
	TestCompleted  int64
	Rejected       int64
}

type reservationRepo interface {
	CreateReservationWithTests(ctx context.Context, userID int, sampleName string, sampleAmount float64, receivedNote *string, sampleResidue string, specialInstructions *string, sampleType string, note *string, batchNumber *string, sampleTakenDate time.Time, sampleTakenBy string, testRequestNumber *string, tests []models.TestItem) (*models.Reservation, error)
	GetAllReservations(ctx context.Context, search string, status []string, page, limit int, sortBy, sortOrder string) ([]models.Reservation, int64, error)
	GetReservationStats(ctx context.Context) (map[string]int64, int64, error)
	UpdateReservationByXid(ctx context.Context, xid uuid.UUID, poURL, poNumber, samplePhoto *string, updatedBy int) (*models.Reservation, error)
	GetReservationByXid(ctx context.Context, xid uuid.UUID) (*models.Reservation, error)
	UpdateReservationStatus(ctx context.Context, xid uuid.UUID, status string, updatedBy int, samplePhoto *string) (*models.Reservation, error)
}

type sequenceStore interface {
	Next(ctx context.Context, key string, ttl time.Duration) (int64, error)
}

type ReservationService struct {
	repo reservationRepo
	seq  sequenceStore
}

var ErrFailedCreate = errors.New("failed to create reservation")
var ErrFailedFetch = errors.New("failed to fetch reservation")
var ErrFailedUpdate = errors.New("failed to update reservation")

func NewReservationService(repo reservationRepo, seq sequenceStore) *ReservationService {
	return &ReservationService{
		repo: repo,
		seq:  seq,
	}
}

func (rs *ReservationService) toReservationResult(ctx context.Context, r models.Reservation) ReservationResult {
	var tests []TestResult
	for _, t := range r.Tests {
		tests = append(tests, TestResult{
			ID:            t.ID,
			XID:           t.XID,
			ReservationID: t.ReservationID,
			TestListId:    t.TestListID,
			Method:        t.Method,
			Criteria:      t.Criteria,
			Status:        t.Status,
			TestUnit:      t.TestUnit,
			TestResult:    t.TestResult,
			TestSummary:   t.TestSummary,
			ApprovedDate:  t.ApprovedDate,
			Notes:         t.Notes,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
			TestList:      t.TestList.Label,
			CategoryName:  t.TestList.Category.Name,
		})
	}

	return ReservationResult{
		ID:                  r.ID,
		Xid:                 r.XID,
		UserID:              r.UserID,
		ReservationDate:     r.ReservationDate,
		Status:              r.Status,
		SampleName:          r.SampleName,
		SampleAmount:        r.SampleAmount,
		ReceivedNote:        r.ReceivedNote,
		TestRequestNumber:   r.TestRequestNumber,
		SampleResidue:       r.SampleResidue,
		SpecialInstructions: r.SpecialInstructions,
		POURL:               r.POURL,
		PONumber:            r.PONumber,
		SampleType:          r.SampleType,
		FeedbackRequired:    r.FeedbackRequired,
		MRApproved:          r.MRApproved,
		HeadApproved:        r.HeadApproved,
		MRApprovalDate:      r.MRApprovalDate,
		HeadApprovalDate:    r.HeadApprovalDate,
		SamplePhoto:         FileURL(FileCategoryReservation, r.SamplePhoto),
		Note:                r.Note,
		StartTestDate:       r.StartTestDate,
		EndTestDate:         r.EndTestDate,
		SAPNotes:            r.SAPNotes,
		BatchNumber:         r.BatchNumber,
		LeadTime:            r.LeadTime,
		SampleTakenDate:     r.SampleTakenDate,
		SampleTakenBy:       r.SampleTakenBy,
		IsKan:               r.IsKan,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
		UpdatedBy:           r.UpdateBy,
		Tests:               tests,
		User: ReservationUser{
			ID:          r.User.ID,
			XID:         r.User.XID,
			Email:       r.User.Email,
			PhoneNumber: r.User.PhoneNumber,
			Name:        r.User.Name,
			Company:     r.User.Company,
			Address:     r.User.Address,
			Avatar:      FileURL(FileCategoryUser, r.User.Avatar),
		},
	}
}

func (rs *ReservationService) generateTRN(ctx context.Context, sampleType string, date time.Time) (string, error) {
	aaa, ok := sampleTypeCode[sampleType]
	if !ok {
		return "", fmt.Errorf("unknown sample type: %s", sampleType)
	}

	dayKey := date.Format("02012006")
	redisKey := fmt.Sprintf("trn:%s", dayKey)

	seq, err := rs.seq.Next(ctx, redisKey, 48*time.Hour)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%03d", aaa, dayKey, seq), nil
}

func (rs *ReservationService) CreateReservation(ctx context.Context, userID int, req CreateReservationParams) (CreateReservationResult, error) {
	trn, err := rs.generateTRN(ctx, req.SampleType, time.Now())
	if err != nil {
		return CreateReservationResult{}, ErrFailedCreate
	}

	var reservation *models.Reservation
	tests := make([]models.TestItem, 0, len(req.Tests))
	for _, t := range req.Tests {
		tests = append(tests, models.TestItem{
			TestListID: t.TestListID,
			Method:     t.Method,
			Criteria:   t.Criteria,
			TestUnit:   t.TestUnit,
		})
	}
	reservation, err = rs.repo.CreateReservationWithTests(ctx,
		userID,
		req.SampleName,
		req.SampleAmount,
		req.ReceivedNote,
		req.SampleResidue,
		req.SpecialInstructions,
		req.SampleType,
		req.Note,
		req.BatchNumber,
		req.SampleTakenDate,
		req.SampleTakenBy,
		&trn,
		tests,
	)

	if err != nil {
		return CreateReservationResult{}, ErrFailedCreate
	}

	return CreateReservationResult{
		ID:                  reservation.ID,
		Xid:                 reservation.XID,
		SampleName:          reservation.SampleName,
		SampleAmount:        reservation.SampleAmount,
		SpecialInstructions: reservation.SpecialInstructions,
		SampleType:          reservation.SampleType,
		BatchNumber:         reservation.BatchNumber,
		SampleTakenDate:     reservation.SampleTakenDate,
		SampleTakenBy:       reservation.SampleTakenBy,
		Message:             "Reservation created",
		TestRequestNumber:   reservation.TestRequestNumber,
		Tests:               req.Tests,
	}, nil
}

func (rs *ReservationService) ListReservations(ctx context.Context, params ListParams) (PaginatedReservationResult, error) {
	reservations, total, err := rs.repo.GetAllReservations(ctx, params.Search, params.Status, params.Page, params.Limit, params.SortBy, params.SortOrder)
	if err != nil {
		return PaginatedReservationResult{}, ErrFailedFetch
	}
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	var response []ReservationResult

	for _, r := range reservations {
		result := rs.toReservationResult(ctx, r)
		result.Message = "Reservations Fetched Successfully"
		response = append(response, result)
	}

	return PaginatedReservationResult{
		Data:       response,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (rs *ReservationService) GetReservationStats(ctx context.Context) (ReservationStatsResult, error) {
	counts, total, err := rs.repo.GetReservationStats(ctx)
	if err != nil {
		return ReservationStatsResult{}, ErrFailedFetch
	}
	return ReservationStatsResult{
		Total:          total,
		Pending:        counts["pending"],
		SampleReceived: counts["sample_received"],
		TestOnProgress: counts["test_on_progress"],
		TestCompleted:  counts["test_completed"],
		Rejected:       counts["rejected"],
	}, nil
}

func (rs *ReservationService) UpdateReservationByXid(ctx context.Context, updatedBy int, xid uuid.UUID, req UpdateReservationParams) (MessageResult, error) {
	_, err := rs.repo.UpdateReservationByXid(ctx, xid, req.POURL, req.PONumber, req.SamplePhoto, updatedBy)
	if err != nil {
		return MessageResult{}, ErrFailedUpdate
	}

	return MessageResult{
		Message: "Reservation updated",
	}, nil
}

func (rs *ReservationService) GetReservationByXid(ctx context.Context, xid uuid.UUID) (ReservationResult, error) {
	reservation, err := rs.repo.GetReservationByXid(ctx, xid)
	if err != nil {
		return ReservationResult{}, ErrFailedFetch
	}

	result := rs.toReservationResult(ctx, *reservation)
	result.Message = "Reservation fetched"
	return result, nil
}

func (rs *ReservationService) UpdateReservationStatus(ctx context.Context, updatedBy int, xid uuid.UUID, req UpdateStatusParams) (MessageResult, error) {
	reservation, err := rs.repo.GetReservationByXid(ctx, xid)
	if err != nil {
		return MessageResult{}, ErrFailedFetch
	}

	currentStatus := ""
	if reservation.Status != nil {
		currentStatus = *reservation.Status
	}

	hasSamplePhoto := req.SamplePhoto != nil
	if err := lifecycle.Next(lifecycle.Status(currentStatus), lifecycle.Status(req.Status), hasSamplePhoto); err != nil {
		return MessageResult{}, err
	}

	_, err = rs.repo.UpdateReservationStatus(ctx, xid, req.Status, updatedBy, req.SamplePhoto)
	if err != nil {
		return MessageResult{}, ErrFailedUpdate
	}

	return MessageResult{
		Message: "Reservation status updated",
	}, nil
}
