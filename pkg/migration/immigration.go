package migration

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
)

type immigration struct {
	ctx        context.Context
	migrations *source.Migrations
}

func (i immigration) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (i immigration) Close() error {
	return nil
}

func (i immigration) First() (version uint, err error) {
	version, ok := i.migrations.First()
	if !ok {
		return 0, &os.PathError{
			Op:   "first",
			Path: "migrate",
			Err:  os.ErrNotExist,
		}
	}

	return version, nil
}

func (i immigration) Prev(version uint) (prevVersion uint, err error) {
	prevVersion, ok := i.migrations.Prev(version)
	if !ok {
		return 0, &os.PathError{
			Op:   fmt.Sprintf("prev for version %v", version),
			Path: "migrate",
			Err:  os.ErrNotExist,
		}
	}

	return prevVersion, nil
}

func (i immigration) Next(version uint) (nextVersion uint, err error) {
	nextVersion, ok := i.migrations.Next(version)
	if !ok {
		return 0, &os.PathError{
			Op:   fmt.Sprintf("prev for version %v", version),
			Path: "migrate",
			Err:  os.ErrNotExist,
		}
	}

	return nextVersion, nil
}

func (i immigration) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	m, ok := i.migrations.Up(version)
	if !ok {
		return nil, "", &os.PathError{
			Op:   fmt.Sprintf("read version %v", version),
			Path: "migrate",
			Err:  os.ErrNotExist,
		}
	}

	return ioutil.NopCloser(strings.NewReader(m.Raw)), m.Identifier, nil
}

func (i immigration) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	m, ok := i.migrations.Down(version)
	if !ok {
		return nil, "", &os.PathError{
			Op:   fmt.Sprintf("read version %v", version),
			Path: "migrate",
			Err:  os.ErrNotExist,
		}
	}

	return ioutil.NopCloser(strings.NewReader(m.Raw)), m.Identifier, nil
}

func Immigration(ctx context.Context, tenantId string, migrations []Migrate) (source.Driver, error) {
	mig := source.NewMigrations()

	for i, m := range migrations {
		v := uint(i + 1)
		sqlUp, err := m.Up(ctx, tenantId)
		if err != nil {
			return nil, err
		}

		sqlDown, err := m.Down(ctx, tenantId)
		if err != nil {
			return nil, err
		}

		mig.Append(&source.Migration{
			Version:    v,
			Identifier: m.ID(ctx),
			Direction:  source.Up,
			Raw:        sqlUp,
		})

		mig.Append(&source.Migration{
			Version:    v,
			Identifier: m.ID(ctx),
			Direction:  source.Down,
			Raw:        sqlDown,
		})
	}

	return &immigration{
		migrations: mig,
	}, nil
}
