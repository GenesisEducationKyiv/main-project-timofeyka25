package storage

import (
	"genesis-test/src/app/service"
	"genesis-test/src/app/utils"
)

type csvStorage struct {
	filepath string
}

func NewCsvStorage(filepath string) service.EmailStorage {
	return &csvStorage{filepath: filepath}
}

func (c *csvStorage) GetAllEmails() ([]string, error) {
	subscribed, err := utils.ReadAllFromCsvToSlice(c.filepath)
	if err != nil {
		return nil, err
	}

	return subscribed, nil
}

func (c *csvStorage) AddEmail(newEmail string) error {
	emails, err := c.GetAllEmails()
	if err != nil {
		return err
	}
	sorted, err := utils.InsertToSortedSlice(emails, newEmail)
	if err != nil {
		return err
	}
	return utils.WriteToCsv(c.filepath, sorted)
}
