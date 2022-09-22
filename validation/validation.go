package validation

import (
	"errors"
	"net/http"
)

func Validation(link string) error {
	response, err := http.Get(link)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Bad link!")
	}
	return nil
}
