package db

import (
	"io"

	"github.com/pkg/errors"
)

// NewConnectionGoPG will create new connection for selected package.
func NewConnectionGoPG(config Config) (sql SQL, err error) {
	writerConf := config.Master
	writerConn, writerCloser, err := connectorGoPgWriter(writerConf, config.Logger)
	if err != nil {
		return
	}

	var slaves = make([]slaveMap, 0)
	for _, slaveConf := range config.Slaves {
		readerConn, readerCloser, readerErr := connectorGoPgWriter(slaveConf, config.Logger)
		if readerErr != nil {
			err = readerErr
			return
		}

		slaves = append(slaves, slaveMap{
			conf:   slaveConf,
			conn:   readerConn,
			closer: readerCloser,
		})
	}

	return &goPgSQL{
		masterConf:   writerConf,
		masterConn:   writerConn,
		masterCloser: writerCloser,
		slaves:       slaves,
	}, nil
}

// slaveMap configuration and connection mapping for slave
type slaveMap struct {
	conf   Conf
	conn   SQLReader
	closer io.Closer
}

// goPgSQL is a struct implements SQL interface
type goPgSQL struct {
	masterConf   Conf
	masterConn   SQLWriter
	masterCloser io.Closer
	slaves       []slaveMap
}

// Writer always use the master.
func (g *goPgSQL) Writer() SQLWriter {
	return g.masterConn
}

// Reader using slaves instance.
// Fallback using master
// TODO: selecting slave
func (g *goPgSQL) Reader() SQLReader {
	for _, slave := range g.slaves {
		if slave.conn != nil {
			return slave.conn
		}
	}

	return g.masterConn
}

func (g *goPgSQL) Close() error {
	var err error
	errMaster := g.masterCloser.Close()
	if errMaster != nil {
		err = errors.Wrapf(err, errMaster.Error())
	}

	for _, c := range g.slaves {
		errReader := c.closer.Close()
		if errReader != nil {
			err = errors.Wrapf(err, errReader.Error())
		}
	}

	return err
}
