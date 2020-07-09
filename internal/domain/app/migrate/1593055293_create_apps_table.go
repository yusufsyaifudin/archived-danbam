package migrate

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
)

// CreateAppsTable1593055293 is struct to define a migration with ID 1593055293_create_apps_table
type CreateAppsTable1593055293 struct{}

// ID return unique identifier for each migration. The prefix is unix time when this migration is created.
func (m CreateAppsTable1593055293) ID(ctx context.Context) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateAppsTable1593055293.ID")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return fmt.Sprintf("%d_%s.sql", 1593055293, "create_apps_table")
}

// SequenceNumber return current time when the migration is created,
// this useful to see the current status of the migration.
func (m CreateAppsTable1593055293) SequenceNumber(ctx context.Context) int {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateAppsTable1593055293.SequenceNumber")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	return 1593055293
}

// Up return sql migration for sync database
func (m CreateAppsTable1593055293) Up(ctx context.Context, tenantID string) (sql string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateAppsTable1593055293.Up")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	sql = `
CREATE TABLE IF NOT EXISTS apps (
	id VARCHAR NOT NULL PRIMARY KEY,
	name VARCHAR NOT NULL DEFAULT '',
	enabled BOOL NOT NULL DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
`

	return
}

// Down return sql migration for rollback database
func (m CreateAppsTable1593055293) Down(ctx context.Context, tenantID string) (sql string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateAppsTable1593055293.Down")
	defer func() {
		ctx.Done()
		span.Finish()
	}()

	sql = `DROP TABLE IF EXISTS apps;`
	return
}
