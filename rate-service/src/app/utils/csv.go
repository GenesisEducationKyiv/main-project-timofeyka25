package utils

import (
	"encoding/csv"
	"genesis-test/src/app/customerror"
	"os"

	"github.com/pkg/errors"
)

func WriteToCsv(path string, data []string) (err error) {
	if path == "" || len(data) < 1 {
		return customerror.ErrInvalidInput
	}
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer func(f *os.File) {
		err = f.Close()
	}(file)
	w := csv.NewWriter(file)

	for _, v := range data {
		err = w.Write([]string{v})
		if err != nil {
			return errors.Wrap(err, "write to csv")
		}
	}

	defer w.Flush()
	err = w.Error()
	return err
}

func ReadAllFromCsvToSlice(path string) (data []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, customerror.ErrFileNotFound
	}
	defer func(f *os.File) {
		err = f.Close()
	}(file)

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "read all from csv")
	}

	for _, record := range records {
		data = append(data, record[0])
	}

	return data, nil
}
