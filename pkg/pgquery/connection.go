package pgquery

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.com/aosorgin/pggen/pkg/datagen"
)

type pgConnection struct {
	conn *pgx.Conn

	query string
	data  datagen.DataGenerator
}

func (c *pgConnection) Query(ctx context.Context) error {
	var args []interface{}
	var err error
	if c.data != nil {
		strArgs, err := c.data.Next()
		if err != nil {
			return errors.Wrap(err, "failed to generate data")
		}

		args = make([]interface{}, len(strArgs))
		for i, val := range strArgs {
			args[i] = val
		}
	}
	_, err = c.conn.Query(ctx, c.query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	return nil
}

func (c *pgConnection) Close(ctx context.Context) error {
	if err := c.conn.Close(ctx); err != nil {
		return errors.Wrap(err, "failed to close connection")
	}

	return nil
}

func NewPgConnection(ctx context.Context, dbConn, query string, data datagen.DataGenerator) (*pgConnection, error) {
	conn, err := pgx.Connect(context.Background(), dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return &pgConnection{
		conn:  conn,
		query: query,
		data:  data,
	}, nil
}
