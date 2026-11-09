package repo

import (
	"context"
	"labflow/models"
	"labflow/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestRepo struct {
	q *sqlc.Queries
}

func NewTestRepo(pool *pgxpool.Pool) *TestRepo {
	return &TestRepo{q: sqlc.New(pool)}
}

func (t *TestRepo) GetTestByXid(ctx context.Context, xid uuid.UUID) (*models.Test, error) {
	row, err := t.q.GetTestWithCategory(ctx, xid)
	if err != nil {
		return nil, wrapDBError(err)
	}
	test := toTest(row.Test)
	test.TestList = toTestListWithCategory(row.TestList, row.Category)
	return &test, nil
}

func (t *TestRepo) UpdateTestByXid(ctx context.Context, xid uuid.UUID, testResult, testSummary, notes *string, updatedBy int) (*models.Test, error) {
	row, err := t.q.UpdateTestResult(ctx, sqlc.UpdateTestResultParams{
		Xid:         xid,
		TestResult:  testResult,
		TestSummary: testSummary,
		Notes:       notes,
		UpdatedBy:   &updatedBy,
	})
	if err != nil {
		return nil, wrapDBError(err)
	}
	test := toTest(row)
	return &test, nil
}

func (t *TestRepo) UpdateIsApproved(ctx context.Context, xid uuid.UUID, isApproved *bool, updatedBy int) (*models.Test, error) {
	row, err := t.q.UpdateTestApproval(ctx, sqlc.UpdateTestApprovalParams{
		Xid:        xid,
		IsApproved: isApproved,
		UpdatedBy:  &updatedBy,
	})
	if err != nil {
		return nil, wrapDBError(err)
	}
	test := toTest(row)
	return &test, nil
}
