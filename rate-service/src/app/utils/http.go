package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetAndParse(url string, target interface{}) (err error) {
	r, err := http.Get(url) //nolint:noctx
	if err != nil {
		return err
	}
	defer func() {
		err = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}
