package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"labflow/sqlc"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepo struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func NewReservationRepo(pool *pgxpool.Pool) *ReservationRepo {
	return &ReservationRepo{pool: pool, q: sqlc.New(pool)}
}

func (r *ReservationRepo) CreateReservationWithTests(ctx context.Context, userID int, sampleName string, sampleAmount float64, receivedNote *string, sampleResidue string, specialInstructions *string, sampleType string, note *string, batchNumber *string, sampleTakenDate time.Time, sampleTakenBy string, testRequestNumber *string, tests []models.TestItem) (*models.Reservation, error) {
	var reservation models.Reservation
	err := inTx(ctx, r.pool, func(q *sqlc.Queries) error {
		row, err := q.CreateReservation(ctx, sqlc.CreateReservationParams{
			Xid:                 uuid.New(),
			UserID:              userID,
			ReservationDate:     time.Now(),
			SampleName:          sampleName,
			SampleAmount:        sampleAmount,
			ReceivedNote:        receivedNote,
			TestRequestNumber:   testRequestNumber,
			SampleResidue:       sampleResidue,
			SpecialInstructions: specialInstructions,
			SampleType:          sampleType,
			Note:                note,
			BatchNumber:         batchNumber,
			SampleTakenDate:     sampleTakenDate,
			SampleTakenBy:       sampleTakenBy,
		})
		if err != nil {
			return fmt.Errorf("failed to create reservation: %w", wrapDBError(err))
		}
		reservation = toReservation(row)

		if len(tests) == 0 {
			return nil
		}
		params := sqlc.CreateTestsParams{ReservationID: reservation.ID}
		for _, t := range tests {
			params.Xids = append(params.Xids, uuid.New())
			params.TestListIds = append(params.TestListIds, t.TestListID)
			params.Methods = append(params.Methods, t.Method)
			params.Criteria = append(params.Criteria, t.Criteria)
			params.TestUnits = append(params.TestUnits, t.TestUnit)
		}
		rows, err := q.CreateTests(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to create tests: %w", wrapDBError(err))
		}
		for _, t := range rows {
			reservation.Tests = append(reservation.Tests, toTest(t))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

var allowedReservationSort = map[string]bool{
	"test_request_number": true,
	"sample_name":         true,
	"status":              true,
	"reservation_date":    true,
}

func (r *ReservationRepo) GetAllReservations(ctx context.Context, search string, status []string, page, limit int, sortBy, sortOrder string) ([]models.Reservation, int64, error) {
	pattern := likePattern(search)
	if status == nil {
		status = []string{}
	}

	total, err := r.q.CountReservations(ctx, sqlc.CountReservationsParams{Search: pattern, Statuses: status})
	if err != nil {
		return nil, 0, wrapDBError(err)
	}

	column := sortBy
	if !allowedReservationSort[column] {
		column = "created_at"
	}
	direction := "desc"
	if strings.ToLower(sortOrder) == "asc" {
		direction = "asc"
	}
	pageLimit, offset := pageBounds(page, limit)

	rows, err := r.q.ListReservations(ctx, sqlc.ListReservationsParams{
		Search:     pattern,
		Statuses:   status,
		Sort:       column + ":" + direction,
		PageLimit:  pageLimit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, 0, wrapDBError(err)
	}

	reservations := make([]models.Reservation, len(rows))
	ids := make([]int, len(rows))
	index := make(map[int]int, len(rows))
	for i, row := range rows {
		reservations[i] = toReservation(row.Reservation)
		reservations[i].User = toUser(row.User)
		ids[i] = row.Reservation.ID
		index[row.Reservation.ID] = i
	}
	if len(ids) == 0 {
		return reservations, total, nil
	}

	tests, err := r.q.ListTestsByReservations(ctx, ids)
	if err != nil {
		return nil, 0, wrapDBError(err)
	}
	for _, row := range tests {
		test := toTest(row.Test)
		test.TestList = toTestList(row.TestList)
		res := &reservations[index[row.Test.ReservationID]]
		res.Tests = append(res.Tests, test)
	}
	return reservations, total, nil
}

func (r *ReservationRepo) GetReservationStats(ctx context.Context) (map[string]int64, int64, error) {
	rows, err := r.q.CountReservationsByStatus(ctx)
	if err != nil {
		return nil, 0, wrapDBError(err)
	}

	counts := make(map[string]int64)
	var total int64
	for _, row := range rows {
		counts[deref(row.Status)] = row.Count
		total += row.Count
	}
	return counts, total, nil
}

func (r *ReservationRepo) UpdateReservationByXid(ctx context.Context, xid uuid.UUID, poURL, poNumber, samplePhoto *string, updatedBy int) (*models.Reservation, error) {
	row, err := r.q.UpdateReservationDocuments(ctx, sqlc.UpdateReservationDocumentsParams{
		Xid:         xid,
		PoUrl:       poURL,
		PoNumber:    poNumber,
		SamplePhoto: samplePhoto,
		UpdatedBy:   &updatedBy,
	})
	if err != nil {
		return nil, wrapDBError(err)
	}
	updated := toReservation(row)
	return &updated, nil
}

// GetReservationByXid returns the reservation with its user, and its tests
// with each test's catalog entry and category.
func (r *ReservationRepo) GetReservationByXid(ctx context.Context, xid uuid.UUID) (*models.Reservation, error) {
	row, err := r.q.GetReservationWithUser(ctx, xid)
	if err != nil {
		return nil, wrapDBError(err)
	}
	reservation := toReservation(row.Reservation)
	reservation.User = toUser(row.User)

	tests, err := r.q.ListTestsWithCategory(ctx, reservation.ID)
	if err != nil {
		return nil, wrapDBError(err)
	}
	for _, t := range tests {
		test := toTest(t.Test)
		test.TestList = toTestListWithCategory(t.TestList, t.Category)
		reservation.Tests = append(reservation.Tests, test)
	}
	return &reservation, nil
}

func (r *ReservationRepo) UpdateReservationStatus(ctx context.Context, xid uuid.UUID, status string, updatedBy int, samplePhoto *string) (*models.Reservation, error) {
	row, err := r.q.UpdateReservationStatus(ctx, sqlc.UpdateReservationStatusParams{
		Xid:         xid,
		Status:      &status,
		UpdatedBy:   &updatedBy,
		SamplePhoto: samplePhoto,
	})
	if err != nil {
		return nil, wrapDBError(err)
	}
	updated := toReservation(row)
	return &updated, nil
}
