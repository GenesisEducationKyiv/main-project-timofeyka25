package utils

import (
	"genesis-test/src/app/errors"
	"sort"
)

func InsertToSortedSlice(s []string, toInsert string) ([]string, error) {
	index := sort.SearchStrings(s, toInsert)
	if index < len(s) && s[index] == toInsert {
		return nil, errors.ErrAlreadyExists
	}

	s = append(s[:index], append([]string{toInsert}, s[index:]...)...)

	return s, nil
}
