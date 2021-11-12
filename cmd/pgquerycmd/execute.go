package pgquerycmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gitlab.com/aosorgin/pggen/pkg/datagen"
	"gitlab.com/aosorgin/pggen/pkg/pgquery"
)

var executeCmd = &cobra.Command{
	Use:   "execute <query>",
	Short: "Execute query",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runQueryExecutor,
}

var (
	csvFilePath string
	connCount   int
	queryCount  int
	queryDelay  time.Duration
)

func runQueryExecutor(cmd *cobra.Command, args []string) error {
	query := args[0]

	var err error
	var data datagen.DataGenerator
	if csvFilePath != "" {
		if data, err = datagen.NewCSVFile(csvFilePath); err != nil {
			panic(errors.Wrap(err, "failed to open csv file with query's parameters"))
		}
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	executor := pgquery.StartPgExecutor(ctx, dbConn, query, data, connCount, queryCount, queryDelay)
	executor.Wait()

	return nil
}
