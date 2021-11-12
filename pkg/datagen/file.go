package datagen

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/pkg/errors"
)

type csvDataFile struct {
	rawFile   io.ReadCloser
	cvsReader *csv.Reader
}

func (df *csvDataFile) Next() ([]string, error) {
	return df.cvsReader.Read()
}

func (df *csvDataFile) Close() error {
	df.cvsReader = nil
	return df.rawFile.Close()
}

func NewCSVFile(filePath string) (DataGenerator, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open csv file")
	}

	return &csvDataFile{
		rawFile:   f,
		cvsReader: csv.NewReader(f),
	}, nil
}
