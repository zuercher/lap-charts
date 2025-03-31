package csvx

import (
	"encoding/csv"
	"io"
	"iter"
	"os"
)

func New(file string) (*CSV, error) {
	handle, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(handle)
	reader.FieldsPerRecord = -1

	return &CSV{handle, reader}, nil
}

type CSV struct {
	handle *os.File
	reader *csv.Reader
}

func (c*CSV) Close() error {
	c.reader = nil
	return c.handle.Close()
}

func (c *CSV) Rows() iter.Seq2[[]string, error] {
	return func(yield func([]string, error) bool) {
		defer c.Close()
		for {
			record, err := c.reader.Read()
			if err != nil {
				if err != io.EOF {
					yield(nil, err)
				}
				break
			}

			if !yield(record, nil) {
				break
			}
		}
	}
}
