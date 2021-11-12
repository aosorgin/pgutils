package pgquery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.com/aosorgin/pggen/pkg/datagen"
)

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

func StartPgExecutor(ctx context.Context, dbConn, query string, data datagen.DataGenerator, connCount, queryCount int, queryDelay time.Duration) *PgQueryExecutor {
	res := &PgQueryExecutor{}

	res.ctx, res.cancel = context.WithCancel(ctx)

	for i := 0; i < connCount; i++ {
		go func() {
			conn, err := NewPgConnection(res.ctx, dbConn, query, data)
			if err != nil {
				log.Error(err.Error())
			}

			ticker := time.After(time.Second)
			for q := 0; q < queryCount; q++ {
				select {
				case <-ticker:
					if err := conn.Query(res.ctx); err != nil {
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