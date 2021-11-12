package datagen

import "io"

type DataGenerator interface {
	io.Closer

	Next() ([]string, error)
}
