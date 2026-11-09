package repo

import (
	"context"
	"errors"
	"labflow/apperr"
	"labflow/models"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// These tests run against a migrated and seeded Postgres (`main seed`).
// Set TEST_DATABASE_URL to run them; CI and loadtest/bench provide one.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping repository integration test")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func ptr[T any](v T) *T { return &v }

func unique(prefix string) string { return prefix + "-" + uuid.NewString()[:8] }

func TestAuthRepoSeededAdminHasAllPermissions(t *testing.T) {
	r := NewAuthRepo(testPool(t))
	ctx := context.Background()

	admin, err := r.GetUserByUsername(ctx, "admin")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	user, err := r.GetUserByID(ctx, admin.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if len(user.Roles) != 1 || deref(user.Roles[0].RoleName) != "admin" {
		t.Fatalf("roles = %+v, want [admin]", user.Roles)
	}
	accesses := user.Roles[0].Accesses
	if len(accesses) != 28 {
		t.Fatalf("admin has %d accesses, want 28", len(accesses))
	}
	for i := 1; i < len(accesses); i++ {
		if accesses[i-1].ID >= accesses[i].ID {
			t.Fatalf("accesses not ordered by id: %d then %d", accesses[i-1].ID, accesses[i].ID)
		}
	}
}

func TestAuthRepoCreateAndUpdateUser(t *testing.T) {
	r := NewAuthRepo(testPool(t))
	ctx := context.Background()
	name := unique("user")
	email := name + "@test.local"

	created, err := r.CreateUser(ctx, name, "hash", email, unique("08"), "Test User", "ACME", "Street 1", RoleCustomer)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created.ID == 0 || created.CreatedAt == nil || deref(created.IsVerified) {
		t.Fatalf("created = %+v; want an ID, created_at and is_verified false", created)
	}

	_, err = r.CreateUser(ctx, unique("user"), "hash", email, unique("08"), "Dup", "", "", RoleCustomer)
	if !errors.Is(err, apperr.ErrDuplicateKey) {
		t.Fatalf("duplicate email: got %v, want ErrDuplicateKey", err)
	}

	withRoles, err := r.GetUserByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if len(withRoles.Roles) != 1 || withRoles.Roles[0].ID != RoleCustomer || len(withRoles.Roles[0].Accesses) == 0 {
		t.Fatalf("roles = %+v, want customer with accesses", withRoles.Roles)
	}

	verified, err := r.UpdateUserVerifiedStatus(ctx, email)
	if err != nil || !deref(verified.IsVerified) {
		t.Fatalf("UpdateUserVerifiedStatus: %v, verified=%v", err, deref(verified.IsVerified))
	}
	updated, err := r.UpdateUserPassword(ctx, email, "new-hash")
	if err != nil || updated.Password != "new-hash" {
		t.Fatalf("UpdateUserPassword: %v, password=%q", err, updated.Password)
	}
	byEmail, err := r.GetUserByEmail(ctx, email)
	if err != nil || byEmail.ID != created.ID {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	if _, err := r.GetUserByID(ctx, -1); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("missing user: got %v, want ErrNotFound", err)
	}
	if _, err := r.UpdateUserVerifiedStatus(ctx, "nobody@test.local"); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("verify missing user: got %v, want ErrNotFound", err)
	}
}

func TestUserRepoUpdateSetsNilFieldsToNull(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	auth := NewAuthRepo(pool)
	name := unique("profile")
	created, err := auth.CreateUser(ctx, name, "hash", name+"@test.local", unique("08"), "P", "ACME", "Street", RoleCustomer)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	updated, err := NewUserRepo(pool).UpdateUserByXid(ctx, created.XID, nil, ptr("New Co"), nil, ptr("a.png"), created.ID)
	if err != nil {
		t.Fatalf("UpdateUserByXid: %v", err)
	}
	if updated.Company != "New Co" || deref(updated.Avatar) != "a.png" || updated.PhoneNumber != "" || updated.Address != "" {
		t.Fatalf("updated = %+v", updated)
	}
	if deref(updated.UpdatedBy) != created.ID {
		t.Fatalf("updated_by = %v, want %d", updated.UpdatedBy, created.ID)
	}
}

func createReservation(t *testing.T, r *ReservationRepo, userID int, name string, batch *string) *models.Reservation {
	t.Helper()
	res, err := r.CreateReservationWithTests(context.Background(), userID, name, 2.5, nil, "Return", ptr("handle cold"),
		"Liquid", nil, batch, time.Now(), "tester", ptr(unique("TRN")), []models.TestItem{
			{TestListID: 1, Method: "USP <61>", Criteria: "< 100 CFU/g", TestUnit: "CFU/g"},
			{TestListID: 11, Method: "USP <791>", Criteria: "6.5 - 7.5", TestUnit: "pH"},
		})
	if err != nil {
		t.Fatalf("CreateReservationWithTests: %v", err)
	}
	return res
}

func TestReservationRepoLifecycle(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	r := NewReservationRepo(pool)
	customer, err := NewAuthRepo(pool).GetUserByUsername(ctx, "customer")
	if err != nil {
		t.Fatalf("customer: %v", err)
	}

	name := unique("sample")
	res := createReservation(t, r, customer.ID, name, nil)
	if res.ID == 0 || len(res.Tests) != 2 || res.Tests[0].ID == 0 || res.Tests[0].ReservationID != res.ID {
		t.Fatalf("created = %+v", res)
	}
	if deref(res.Status) != "pending" || deref(res.BatchNumber) != "-" {
		t.Fatalf("defaults: status=%q batch=%q, want pending and -", deref(res.Status), deref(res.BatchNumber))
	}
	withBatch := createReservation(t, r, customer.ID, name, ptr("B-7"))
	if deref(withBatch.BatchNumber) != "B-7" {
		t.Fatalf("batch = %q, want B-7", deref(withBatch.BatchNumber))
	}

	list, total, err := r.GetAllReservations(ctx, name, nil, 1, 10, "sample_name", "asc")
	if err != nil {
		t.Fatalf("GetAllReservations: %v", err)
	}
	if total != 2 || len(list) != 2 {
		t.Fatalf("search %q: total=%d len=%d, want 2", name, total, len(list))
	}
	for _, item := range list {
		if item.User.ID != customer.ID || len(item.Tests) != 2 || item.Tests[0].TestList.Label == "" {
			t.Fatalf("list item missing user or tests: %+v", item)
		}
		if item.Tests[0].TestList.Category.Name != "" {
			t.Fatalf("list must not load categories, as before")
		}
	}

	page2, _, err := r.GetAllReservations(ctx, name, nil, 2, 1, "", "")
	if err != nil || len(page2) != 1 || page2[0].ID != res.ID {
		t.Fatalf("page 2 of 1 (newest first) should be the first reservation: %v %+v", err, page2)
	}

	none, total, err := r.GetAllReservations(ctx, name, []string{"rejected"}, 1, 10, "", "")
	if err != nil || total != 0 || len(none) != 0 {
		t.Fatalf("status filter: err=%v total=%d len=%d", err, total, len(none))
	}

	detail, err := r.GetReservationByXid(ctx, res.XID)
	if err != nil {
		t.Fatalf("GetReservationByXid: %v", err)
	}
	if detail.User.Username != "customer" || len(detail.Tests) != 2 || detail.Tests[0].TestList.Category.Name == "" {
		t.Fatalf("detail missing user, tests or category: %+v", detail)
	}
	if _, err := r.GetReservationByXid(ctx, uuid.New()); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("missing reservation: got %v, want ErrNotFound", err)
	}

	updated, err := r.UpdateReservationStatus(ctx, res.XID, "sample_received", customer.ID, ptr("photo.jpg"))
	if err != nil || deref(updated.Status) != "sample_received" || deref(updated.SamplePhoto) != "photo.jpg" {
		t.Fatalf("UpdateReservationStatus with photo: %v %+v", err, updated)
	}
	updated, err = r.UpdateReservationStatus(ctx, res.XID, "test_on_progress", customer.ID, nil)
	if err != nil || deref(updated.SamplePhoto) != "photo.jpg" {
		t.Fatalf("a nil photo must keep the old one: %v %q", err, deref(updated.SamplePhoto))
	}

	docs, err := r.UpdateReservationByXid(ctx, res.XID, ptr("po.pdf"), ptr("PO-1"), nil, customer.ID)
	if err != nil || deref(docs.POURL) != "po.pdf" || docs.SamplePhoto != nil || deref(docs.UpdateBy) != customer.ID {
		t.Fatalf("UpdateReservationByXid: %v %+v", err, docs)
	}

	stats, statTotal, err := r.GetReservationStats(ctx)
	if err != nil || statTotal < 2 || stats["pending"] < 1 {
		t.Fatalf("GetReservationStats: %v %v %d", err, stats, statTotal)
	}
}

func TestTestRepoUpdates(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	customer, err := NewAuthRepo(pool).GetUserByUsername(ctx, "customer")
	if err != nil {
		t.Fatalf("customer: %v", err)
	}
	res := createReservation(t, NewReservationRepo(pool), customer.ID, unique("sample"), nil)
	r := NewTestRepo(pool)

	test, err := r.GetTestByXid(ctx, res.Tests[0].XID)
	if err != nil || test.TestList.Category.Name == "" || deref(test.IsApproved) {
		t.Fatalf("GetTestByXid: %v %+v", err, test)
	}
	updated, err := r.UpdateTestByXid(ctx, test.XID, ptr("12 CFU/g"), ptr("pass"), nil, customer.ID)
	if err != nil || deref(updated.TestResult) != "12 CFU/g" || updated.Notes != nil {
		t.Fatalf("UpdateTestByXid: %v %+v", err, updated)
	}
	approved, err := r.UpdateIsApproved(ctx, test.XID, ptr(true), customer.ID)
	if err != nil || !deref(approved.IsApproved) || approved.ApprovedDate == nil {
		t.Fatalf("UpdateIsApproved: %v %+v", err, approved)
	}
	if _, err := r.GetTestByXid(ctx, uuid.New()); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("missing test: got %v, want ErrNotFound", err)
	}
}

func TestTestListRepoVersioning(t *testing.T) {
	r := NewTestListRepo(testPool(t))
	ctx := context.Background()

	catalog, err := r.GetTestListByCategory(ctx, 1)
	if err != nil || len(catalog) == 0 || catalog[0].Category.Name != "Microbiology" {
		t.Fatalf("GetTestListByCategory: %v %+v", err, catalog)
	}

	value := unique("Assay")
	created, err := r.CreateTestList(ctx, models.TestList{CategoryID: 2, Value: value, Label: value, Leadtime: ptr(3)})
	if err != nil {
		t.Fatalf("CreateTestList: %v", err)
	}
	if created.Version != 1 || !deref(created.IsCurrent) || created.Category.Name != "Chemistry" {
		t.Fatalf("created = %+v; want version 1, current, Chemistry", created)
	}

	next := *created
	next.Version = created.Version + 1
	next.Label = value + " v2"
	v2, err := r.CreateNewVersion(ctx, created.XID, next)
	if err != nil {
		t.Fatalf("CreateNewVersion: %v", err)
	}
	if v2.Version != 2 || v2.XID == created.XID || !deref(v2.IsCurrent) {
		t.Fatalf("v2 = %+v", v2)
	}
	old, err := r.GetTestListByXid(ctx, created.XID)
	if err != nil || deref(old.IsCurrent) {
		t.Fatalf("previous version must be retired: %v %+v", err, old)
	}

	items, total, err := r.ListTestLists(ctx, value, ptr(2), 1, 10)
	if err != nil || total != 2 || len(items) != 2 || items[0].XID != v2.XID {
		t.Fatalf("ListTestLists newest first: %v total=%d %+v", err, total, items)
	}

	retired, err := r.DeactivateByXid(ctx, v2.XID)
	if err != nil || deref(retired.IsCurrent) {
		t.Fatalf("DeactivateByXid: %v %+v", err, retired)
	}
}

func TestQuotationRepo(t *testing.T) {
	r := NewQuotationRepo(testPool(t))
	ctx := context.Background()
	company := unique("Company")

	created, err := r.CreateQuotation(ctx, "buyer@test.local", company, nil, ptr("10 samples"))
	if err != nil || deref(created.Status) != "pending" {
		t.Fatalf("CreateQuotation: %v %+v", err, created)
	}
	list, total, err := r.GetQuotations(ctx, company, 1, 10)
	if err != nil || total != 1 || len(list) != 1 || list[0].XID != created.XID {
		t.Fatalf("GetQuotations: %v total=%d", err, total)
	}
	updated, err := r.UpdateQuotationStatus(ctx, created.XID, "  approved ")
	if err != nil || deref(updated.Status) != "approved" {
		t.Fatalf("UpdateQuotationStatus: %v %+v", err, updated)
	}
	if _, err := r.UpdateQuotationStatus(ctx, uuid.New(), "approved"); !errors.Is(err, apperr.ErrNotFound) {
		t.Fatalf("missing quotation: got %v, want ErrNotFound", err)
	}
}
