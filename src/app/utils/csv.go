package utils

import (
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
)

func WriteToCsv(path string, data []string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer func() {
		err = f.Close()
	}()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, v := range data {
		err = w.Write([]string{v})
		if err != nil {
			return errors.Wrap(err, "write to csv")
		}
	}

	return nil
}

func ReadAllFromCsvToSlice(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	var data []string
	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "read all from csv")
	}

	for _, record := range records {
		data = append(data, record[0])
	}

	return data, nil
}
