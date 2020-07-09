package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"github.com/yusufsyaifudin/danbam/pkg/db"
)

type repoPostgres struct {
	migration *migrate.Migrate
	conn      db.SQL
}

func (p repoPostgres) GetAppByID(ctx context.Context, id string) (app *App, err error) {
	panic("implement me")
}

func (p repoPostgres) CreateApp(ctx context.Context, id, name string) (app *App, err error) {
	panic("implement me")
}

func (p repoPostgres) DisableApp(ctx context.Context, id string) (app *App, err error) {
	panic("implement me")
}

func (p repoPostgres) EnableApp(ctx context.Context, id string) (app *App, err error) {
	panic("implement me")
}

// Close do close all opened connection
func (p repoPostgres) Close() error {
	var (
		err                                          error
		dbConnCloseErr                               error
		migrationSourceCloseErr, migrationDbCloseErr error
	)

	if p.conn != nil {
		dbConnCloseErr = p.conn.Close()
	}

	if dbConnCloseErr != nil {
		err = errors.Wrapf(err, dbConnCloseErr.Error())
	}

	if p.migration != nil {
		migrationSourceCloseErr, migrationDbCloseErr = p.migration.Close()
	}

	if migrationSourceCloseErr != nil {
		err = errors.Wrapf(err, migrationSourceCloseErr.Error())
	}

	if migrationDbCloseErr != nil {
		err = errors.Wrapf(err, migrationDbCloseErr.Error())
	}

	return err
}

// DB return repo interface which implements using PgSQL
func DB(logger *zap.Logger) (service Repo, err error) {
	conf, err := readConfig()
	if err != nil {
		return
	}

	// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	// postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]
	pgConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		conf.Postgres.Master.Username,
		conf.Postgres.Master.Password,
		conf.Postgres.Master.Host,
		conf.Postgres.Master.Port,
		conf.Postgres.Master.Db,
	)

	migration, err := Migration(pgConnStr)
	if err != nil {
		return
	}

	if _, isDirty, verErr := migration.Version(); isDirty || verErr == migrate.ErrNilVersion {
		// will return "no change" in log when dirty is false
		err = migration.Up()
	}

	if err != nil {
		return
	}

	pgMasterConf := db.Conf{
		Debug:        true,
		AppName:      "DanBam",
		Host:         conf.Postgres.Master.Host,
		Port:         conf.Postgres.Master.Port,
		Username:     conf.Postgres.Master.Username,
		Password:     conf.Postgres.Master.Password,
		Database:     conf.Postgres.Master.Db,
		PoolSize:     10,
		ReadTimeout:  2000,
		WriteTimeout: 2000,
	}

	if conf.DBType == DbTypePostgres {
		dbConn, err := db.NewConnectionGoPG(db.Config{
			Logger: logger,
			Master: pgMasterConf,
			Slaves: []db.Conf{},
		})

		if err != nil {
			return nil, err
		}

		err = dbConn.Writer().ExecContext(context.Background(), "SELECT 1;")

		return &repoPostgres{
			migration: migration,
			conn:      dbConn,
		}, err
	}

	panic(fmt.Errorf("unknown db_type: %s on app", conf.DBType))
}
