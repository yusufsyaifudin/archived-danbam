package app

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	Migrate "github.com/yusufsyaifudin/danbam/internal/domain/app/migrate"
	"github.com/yusufsyaifudin/danbam/pkg/migration"
)

// driverSourceName this should be unique identifier in the application scope
const driverSourceName = "danbam_app_migration"

// migrations is slices of migrations files.
// use global variables to make it test-able
var migrations = []migration.Migrate{
	new(Migrate.CreateAppsTable1593055293),
}

// Migration returns migration instance and error
func Migration(connStr string) (m *migrate.Migrate, err error) {
	const tenantId = "" // tenantId is empty because we will not separate data for "apps" in the db
	src, err := migration.Immigration(context.Background(), tenantId, migrations)
	if err != nil {
		return
	}

	source.Register(driverSourceName, src)
	m, err = migrate.NewWithSourceInstance(driverSourceName, src, connStr)
	return
}
