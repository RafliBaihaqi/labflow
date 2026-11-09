package repo

import (
	"context"
	"fmt"
	"labflow/models"
	"labflow/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestListRepo struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func NewTestListRepo(pool *pgxpool.Pool) *TestListRepo {
	return &TestListRepo{pool: pool, q: sqlc.New(pool)}
}

func (tl *TestListRepo) GetTestListByCategory(ctx context.Context, categoryID int) ([]models.TestList, error) {
	rows, err := tl.q.ListTestListsByCategory(ctx, categoryID)
	if err != nil {
		return nil, wrapDBError(err)
	}
	testlists := make([]models.TestList, len(rows))
	for i, row := range rows {
		testlists[i] = toTestListWithCategory(row.TestList, row.Category)
	}
	return testlists, nil
}

func (tl *TestListRepo) ListTestLists(ctx context.Context, search string, categoryID *int, page, limit int) ([]models.TestList, int64, error) {
	pattern := likePattern(search)
	total, err := tl.q.CountTestLists(ctx, sqlc.CountTestListsParams{Search: pattern, CategoryID: categoryID})
	if err != nil {
		return nil, 0, wrapDBError(err)
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	rows, err := tl.q.ListTestLists(ctx, sqlc.ListTestListsParams{
		Search:     pattern,
		CategoryID: categoryID,
		PageLimit:  limit,
		PageOffset: (page - 1) * limit,
	})
	if err != nil {
		return nil, 0, wrapDBError(err)
	}
	testlists := make([]models.TestList, len(rows))
	for i, row := range rows {
		testlists[i] = toTestListWithCategory(row.TestList, row.Category)
	}
	return testlists, total, nil
}

func (tl *TestListRepo) CreateTestList(ctx context.Context, input models.TestList) (*models.TestList, error) {
	if input.XID == uuid.Nil {
		input.XID = uuid.New()
	}
	if input.IsCurrent == nil {
		isCurrent := true
		input.IsCurrent = &isCurrent
	}
	id, err := tl.q.CreateTestList(ctx, createTestListParams(input))
	if err != nil {
		return nil, fmt.Errorf("failed to create test list: %w", wrapDBError(err))
	}
	return tl.getByID(ctx, tl.q, id)
}

func (tl *TestListRepo) GetTestListByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	row, err := tl.q.GetTestListByXid(ctx, xid)
	if err != nil {
		return nil, wrapDBError(err)
	}
	testList := toTestListWithCategory(row.TestList, row.Category)
	return &testList, nil
}

// CreateNewVersion retires the version identified by xid and inserts input as
// its successor, in one transaction.
func (tl *TestListRepo) CreateNewVersion(ctx context.Context, xid uuid.UUID, input models.TestList) (*models.TestList, error) {
	var created *models.TestList
	err := inTx(ctx, tl.pool, func(q *sqlc.Queries) error {
		if err := q.DeactivateTestList(ctx, xid); err != nil {
			return fmt.Errorf("failed to deactivate previous version: %w", wrapDBError(err))
		}
		input.XID = uuid.New()
		id, err := q.CreateTestList(ctx, createTestListParams(input))
		if err != nil {
			return fmt.Errorf("failed to create new version: %w", wrapDBError(err))
		}
		created, err = tl.getByID(ctx, q, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (tl *TestListRepo) DeactivateByXid(ctx context.Context, xid uuid.UUID) (*models.TestList, error) {
	if err := tl.q.DeactivateTestList(ctx, xid); err != nil {
		return nil, fmt.Errorf("failed to deactivate test list: %w", wrapDBError(err))
	}
	return tl.GetTestListByXid(ctx, xid)
}

func (tl *TestListRepo) getByID(ctx context.Context, q *sqlc.Queries, id int) (*models.TestList, error) {
	row, err := q.GetTestListByID(ctx, id)
	if err != nil {
		return nil, wrapDBError(err)
	}
	testList := toTestListWithCategory(row.TestList, row.Category)
	return &testList, nil
}

func createTestListParams(input models.TestList) sqlc.CreateTestListParams {
	return sqlc.CreateTestListParams{
		Xid:                input.XID,
		CategoryID:         input.CategoryID,
		Value:              input.Value,
		Label:              input.Label,
		SampleAmountLiquid: input.SampleAmountLiquid,
		SampleAmountSolid:  input.SampleAmountSolid,
		SampleAmountPowder: input.SampleAmountPowder,
		Leadtime:           input.Leadtime,
		Version:            input.Version,
		IsCurrent:          input.IsCurrent,
	}
}
