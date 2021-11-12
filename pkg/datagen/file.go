package datagen

import (
	"encoding/csv"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type csvDataFile struct {
	mutex sync.Mutex

	rawFile   io.ReadCloser
	cvsReader *csv.Reader
}

func (df *csvDataFile) Next() ([]string, error) {
	df.mutex.Lock()
	defer df.mutex.Unlock()
	return df.cvsReader.Read()
}

func (df *csvDataFile) Close() error {
	df.cvsReader = nil
	return df.rawFile.Close()
}

func NewCSVFile(filePath string, seperator rune) (DataGenerator, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open csv file")
	}

	cvsReader := csv.NewReader(f)
	cvsReader.Comma = seperator

	return &csvDataFile{
		rawFile:   f,
		cvsReader: cvsReader,
	}, nil
}
