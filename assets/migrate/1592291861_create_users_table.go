package migrate

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

// CreateUsersTable1592291861 is struct to define a migration with ID 1592291861_create_users_table
type CreateUsersTable1592291861 struct{}

// ID return unique identifier for each migration. The prefix is unix time when this migration is created.
func (m CreateUsersTable1592291861) ID(ctx context.Context) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1592291861.ID")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return fmt.Sprintf("%d_%s.sql", 1592291861, "create_users_table")
}

// SequenceNumber return current time when the migration is created,
// this useful to see the current status of the migration.
func (m CreateUsersTable1592291861) SequenceNumber(ctx context.Context) int {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1592291861.SequenceNumber")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return 1592291861
}

// Up return sql migration for sync database
func (m CreateUsersTable1592291861) Up(ctx context.Context, tenantID string) (sql string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1592291861.Up")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	sql = `
CREATE
`

	return
}

// Down return sql migration for rollback database
func (m CreateUsersTable1592291861) Down(ctx context.Context, tenantID string) (sql string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateUsersTable1592291861.Down")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return
}
