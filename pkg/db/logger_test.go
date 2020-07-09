package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

var ormDB = newGoPGDBTest()
var z, _ = zap.NewDevelopment()

func TestGoPGDBLogger(t *testing.T) {
	convey.Convey("goPGDBLogger", t, func() {
		convey.Convey("debug is not active, so return immediately", func() {
			dbLogger := &goPGDBLogger{
				debug:     false,
				zapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     nil,
				Params:    nil,
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("debug is active, error occurred because no query passed in event query", func() {
			dbLogger := &goPGDBLogger{
				debug:     true,
				zapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     nil, // since this query is nil, so error will be returned in this block: e.FormattedQuery()
				Params:    nil,
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldNotBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldNotBeNil)
		})

		convey.Convey("debug is active, SELECT 1", func() {

			dbLogger := &goPGDBLogger{
				debug:     true,
				zapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     "SELECT 1",
				Params:    nil,
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldBeNil)
		})

		convey.Convey("debug is active, Other query is success with args", func() {
			dbLogger := &goPGDBLogger{
				debug:     true,
				zapLogger: z,
			}

			ctx := context.Background()

			eventQuery := &pg.QueryEvent{
				StartTime: time.Now(),
				DB:        ormDB,
				Model:     nil,
				Query:     "INSERT INTO users(name) VALUES (?) RETURNING *",
				Params:    []interface{}{"user name"},
				Result:    nil,
				Err:       nil,
				Stash:     nil,
			}

			resCtx, err := dbLogger.BeforeQuery(ctx, eventQuery)
			convey.So(resCtx, convey.ShouldNotBeNil)
			convey.So(err, convey.ShouldBeNil)

			err = dbLogger.AfterQuery(ctx, eventQuery)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}
