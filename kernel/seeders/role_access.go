package seeders

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedRoleAccess(ctx context.Context, pool *pgxpool.Pool) {
	// Join table — raw insert is pragmatic for many-to-many links
	// admin gets all, customer limited, qa + lab roles expanded per OLD_APP_REFERENCE sidebar matrix
	query := `INSERT INTO role_accesses (role_id, access_id) VALUES
		-- admin all 28
		(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10), (1, 11), (1, 12), (1, 13), (1, 14), (1, 15), (1, 16), (1, 17),(1, 18),(1, 19),(1, 20),(1, 21),(1, 22),(1, 23),(1, 24),(1, 25),(1, 26),(1, 27),(1, 28),
		-- customer
		(2, 5), (2, 6), (2, 7), (2, 8), (2, 9), (2, 10),(2, 3), (2, 18), (2, 14), (2, 22), (2, 23), (2, 24), (2, 25), (2, 28),
		-- qa
		(3, 6), (3, 7), (3, 8), (3, 10), (3, 11), (3, 12), (3, 13), (3, 14), (3, 15), (3, 16), (3, 17), (3, 18), (3, 19), (3, 20), (3, 21), (3, 22), (3, 23), (3, 25), (3, 26), (3, 27), (3, 28),
		-- micro-analyst (4)
		(4, 6), (4, 10), (4, 11), (4, 2), (4, 18), (4, 14), (4, 28),
		-- chemical-analyst (5)
		(5, 6), (5, 10), (5, 11), (5, 2), (5, 18), (5, 14), (5, 28),
		-- micro-qc (6)
		(6, 6), (6, 7), (6, 10), (6, 11), (6, 17), (6, 13), (6, 14), (6, 15), (6, 18), (6, 25), (6, 28),
		-- chemical-qc (7)
		(7, 6), (7, 7), (7, 10), (7, 11), (7, 17), (7, 13), (7, 14), (7, 15), (7, 18), (7, 25), (7, 28),
		-- lab-head (8)
		(8, 6), (8, 7), (8, 10), (8, 14), (8, 15), (8, 17), (8, 2), (8, 18), (8, 19), (8, 20), (8, 21), (8, 25), (8, 26), (8, 27), (8, 28),
		-- mr (9)
		(9, 6), (9, 10), (9, 14), (9, 17), (9, 2), (9, 18), (9, 25), (9, 28)
		ON CONFLICT DO NOTHING;`
	if _, err := pool.Exec(ctx, query); err != nil {
		log.Fatalf("Failed to seed role_access data: %v", err)
	}
}
