package pgquerycmd

import (
	"math"
	"time"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "pgquery",
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&dbConn, "db", "", "", "Database connection string")

	executeCmd.Flags().IntVar(&connCount, "conn-count", 1, "connections count")
	executeCmd.Flags().IntVar(&queryCount, "query-count", math.MaxInt, "queries count")
	executeCmd.Flags().DurationVar(&queryDelay, "query-delay", time.Second, "query's delay")
	executeCmd.Flags().StringVar(&csvFilePath, "csv-data", "", "CVS file with query data")

	RootCmd.AddCommand(executeCmd)
}
