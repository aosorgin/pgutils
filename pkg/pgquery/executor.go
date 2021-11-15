package pgquery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aosorgin/pgutils/pkg/datagen"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type QueryProcessor interface {
	Query(ctx context.Context) error
	Close(ctx context.Context) error
}
type PgQueryExecutor struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func (qe *PgQueryExecutor) Wait() {
	qe.wg.Wait()
}

func (qe *PgQueryExecutor) Close() {
	qe.cancel()
	qe.wg.Wait()
}

func StartPgExecutor(ctx context.Context, usePsql bool, dbConn, query string, data datagen.DataGenerator,
	connCount, queryCount int, queryDelay time.Duration) *PgQueryExecutor {
	res := &PgQueryExecutor{}

	res.ctx, res.cancel = context.WithCancel(ctx)

	for i := 0; i < connCount; i++ {
		go func() {
			var proc QueryProcessor
			var err error
			if usePsql {
				proc, err = NewPSQLExecutor(res.ctx, dbConn, query, data)
			} else {
				proc, err = NewPgConnection(res.ctx, dbConn, query, data)
			}
			if err != nil {
				log.Error(err.Error())
			}

			ticker := time.After(time.Second)
			for q := 0; q < queryCount; q++ {
				select {
				case <-ticker:
					if err := proc.Query(res.ctx); err != nil {
						log.Error(errors.Wrap(err, "failed to execute query"))
					} else {
						fmt.Print(".")
					}
					ticker = time.After(queryDelay)
				case <-res.ctx.Done():
					res.wg.Done()
					return
				}
			}
			res.wg.Done()
		}()

		res.wg.Add(1)
	}

	return res
}
