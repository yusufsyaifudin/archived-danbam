package app

import "context"

type Repo interface {
	GetAppByID(ctx context.Context, id string) (app *App, err error)
	CreateApp(ctx context.Context, id, name string) (app *App, err error)
	DisableApp(ctx context.Context, id string) (app *App, err error)
	EnableApp(ctx context.Context, id string) (app *App, err error)

	// closing the connection
	Close() error
}
