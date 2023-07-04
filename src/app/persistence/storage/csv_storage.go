package storage

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/service"
	"genesis-test/src/app/utils"

	"github.com/pkg/errors"
)

type csvRepository struct {
	filepath string
}

func NewCsvRepository(filepath string) service.EmailStorage {
	return &csvRepository{filepath: filepath}
}

func (c *csvRepository) GetAllEmails() ([]string, error) {
	subscribed, err := utils.ReadAllFromCsvToSlice(c.filepath)
	if len(subscribed) < 1 {
		return nil, customerror.ErrNoSubscribers
	}
	if err != nil {
		return nil, err
	}

	return subscribed, nil
}

func (c *csvRepository) AddEmail(newEmail string) error {
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
