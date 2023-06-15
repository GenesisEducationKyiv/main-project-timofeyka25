package utils

import (
	"encoding/csv"
	"os"
)

func WriteToCsv(path string, data []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	w := csv.NewWriter(f)
	defer w.Flush()

	for _, v := range data {
		err = w.Write([]string{v})
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadAllFromCsvToSlice(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	var data []string
	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		data = append(data, record[0])
	}

	return data, nil
}
