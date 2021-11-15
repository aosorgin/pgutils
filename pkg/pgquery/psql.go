package pgquery

import (
	"context"
	"net/url"
	"os/exec"
	"strings"

	"github.com/aosorgin/pgutils/pkg/datagen"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type psqlExecutor struct {
	argv []string

	query string
	data  datagen.DataGenerator
}

func (e *psqlExecutor) updateQueryArg() error {
	if e.data != nil {
		values, err := e.data.Next()
		if err != nil {
			return errors.Wrap(err, "failed to get next values from generator")
		}
		query := strings.Replace(e.query, "?", strings.Join(values, ","), 1)
		e.argv[len(e.argv)-1] = query
	}

	return nil
}

func (e *psqlExecutor) Query(ctx context.Context) error {
	if err := e.updateQueryArg(); err != nil {
		return errors.Wrap(err, "failed to update psql arguments")
	}

	cmd := exec.CommandContext(ctx, "psql", e.argv...)
	log.Debug("run: ", cmd.String())
	if err := cmd.Run(); err != nil {
		output, err := cmd.Output()
		if err != nil {
			log.Error(errors.Wrap(err, "failed to collect command output"))
		}
		log.Error(string(output))
		return errors.Wrap(err, "failed to execute psql utility")
	}

	return nil
}

func (e *psqlExecutor) Close(ctx context.Context) error {
	return nil
}

func NewPSQLExecutor(ctx context.Context, dbConn, query string, data datagen.DataGenerator) (QueryProcessor, error) {
	url, err := url.Parse(dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse database connection string")
	}

	var argv []string
	if url.Hostname() != "" {
		argv = append(argv, "-h")
		argv = append(argv, url.Hostname())
	}

	if url.Port() != "" {
		argv = append(argv, "-p")
		argv = append(argv, url.Port())
	}

	if url.User.Username() != "" {
		argv = append(argv, "-U")
		argv = append(argv, url.User.Username())
	}

	argv = append(argv, "-d")
	argv = append(argv, url.Path[1:]) // Database name
	argv = append(argv, "-c")
	argv = append(argv, query)

	return &psqlExecutor{
		argv:  argv,
		query: query,
		data:  data,
	}, nil
}
