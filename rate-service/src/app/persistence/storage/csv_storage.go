package storage

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/utils"

	"github.com/pkg/errors"
)

type EmailStorage struct {
	filepath string
}

func NewCsvRepository(filepath string) *EmailStorage {
	return &EmailStorage{filepath: filepath}
}

func (c *EmailStorage) GetAllEmails() ([]string, error) {
	subscribed, err := utils.ReadAllFromCsvToSlice(c.filepath)
	if len(subscribed) < 1 {
		return nil, customerror.ErrNoSubscribers
	}
	if err != nil {
		return nil, err
	}

	return subscribed, nil
}

func (c *EmailStorage) AddEmail(newEmail string) error {
	emails, err := c.GetAllEmails()
	if err != nil && !errors.Is(err, customerror.ErrNoSubscribers) {
		return err
	}
	sorted, err := utils.InsertToSortedSlice(emails, newEmail)
	if err != nil {
		return err
	}
	return utils.WriteToCsv(c.filepath, sorted)
}
