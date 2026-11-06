package service

import (
	"context"
	"sync"
	"time"

	"labflow/models"
)

type accessRepo interface {
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
}

// DefaultAccessCacheTTL bounds how long a permission change takes to apply.
// Roles are only changed by seeding/admin SQL today, so 30s is plenty fresh.
const DefaultAccessCacheTTL = 30 * time.Second

type accessEntry struct {
	perms   map[string]struct{}
	expires time.Time
}

// AccessService answers "may user X do Y?". Every protected request asks, so
// each user's permission set is cached in-process for a short TTL instead of
// re-loading the user, roles and accesses from Postgres on every call.
type AccessService struct {
	repo accessRepo
	ttl  time.Duration
	now  func() time.Time

	mu    sync.RWMutex
	cache map[int]accessEntry
}

func NewAccessService(repo accessRepo) *AccessService {
	return NewAccessServiceWithTTL(repo, DefaultAccessCacheTTL)
}

func NewAccessServiceWithTTL(repo accessRepo, ttl time.Duration) *AccessService {
	return &AccessService{repo: repo, ttl: ttl, now: time.Now, cache: map[int]accessEntry{}}
}

// HasAccess reports whether the user holds the named permission. A non-nil
// error means the answer is unknown (e.g. the database is unavailable) and
// must not be treated as a denial.
func (s *AccessService) HasAccess(ctx context.Context, userID int, accessName string) (bool, error) {
	perms, err := s.permissions(ctx, userID)
	if err != nil {
		return false, err
	}
	_, ok := perms[accessName]
	return ok, nil
}

// Invalidate drops a user's cached permissions, e.g. after a role change.
func (s *AccessService) Invalidate(userID int) {
	s.mu.Lock()
	delete(s.cache, userID)
	s.mu.Unlock()
}

func (s *AccessService) permissions(ctx context.Context, userID int) (map[string]struct{}, error) {
	now := s.now()
	s.mu.RLock()
	entry, ok := s.cache[userID]
	s.mu.RUnlock()
	if ok && now.Before(entry.expires) {
		return entry.perms, nil
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	perms := map[string]struct{}{}
	for _, role := range user.Roles {
		for _, access := range role.Accesses {
			if access.AccessName != nil {
				perms[*access.AccessName] = struct{}{}
			}
		}
	}

	s.mu.Lock()
	s.cache[userID] = accessEntry{perms: perms, expires: now.Add(s.ttl)}
	s.mu.Unlock()
	return perms, nil
}
