package db

import (
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/go-pg/pg/v9"
)

type closer struct {
	conn *pg.DB
}

func (c closer) Close() error {
	return c.conn.Close()
}

func newCloser(conn *pg.DB) io.Closer {
	return &closer{
		conn: conn,
	}
}

// connectorGoPgWriter using go-pg connection
func connectorGoPgWriter(conf Conf, logger *zap.Logger) (SQLWriter, io.Closer, error) {
	ormPgDB := pg.Connect(&pg.Options{
		Addr:               fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		User:               conf.Username,
		Password:           conf.Password,
		Database:           conf.Database,
		ApplicationName:    conf.AppName,
		ReadTimeout:        time.Duration(conf.ReadTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(conf.WriteTimeout) * time.Millisecond,
		PoolSize:           conf.PoolSize,
		MinIdleConns:       10,
		MaxConnAge:         5 * time.Minute,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
	})

	ormPgDB.AddQueryHook(NewLogger(conf.Debug, logger))

	writer, err := newGoPgWriter(conf, ormPgDB)
	if err != nil {
		return nil, nil, err
	}

	return writer, newCloser(ormPgDB), nil
}
