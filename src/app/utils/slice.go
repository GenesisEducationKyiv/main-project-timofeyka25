package utils

import (
	"genesis-test/src/app/customerror"
	"sort"
)

func InsertToSortedSlice(s []string, toInsert string) ([]string, error) {
	index := sort.SearchStrings(s, toInsert)
	if index < len(s) && s[index] == toInsert {
		return nil, customerror.ErrAlreadyExists
	}

	s = append(s[:index], append([]string{toInsert}, s[index:]...)...)

	return s, nil
}
