package migration

import "context"

type Migrate interface {
	// ID return unique identifier for each migration. The prefix must be number
	ID(ctx context.Context) string

	// SequenceNumber must be unique, this useful to see the current status of the migration.
	SequenceNumber(ctx context.Context) int

	// Up return sql migration for sync database
	Up(ctx context.Context, tenantID string) (sql string, err error)

	// Down return sql migration for rollback database
	Down(ctx context.Context, tenantID string) (sql string, err error)
}
