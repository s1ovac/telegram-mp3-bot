package validation

import (
	"errors"
	"net/http"
	"net/url"
)

// Validation checks link for spell mistakes
func Validation(link string) error {
	response, err := http.Get(link)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Connection is bad...")
	}
	return nil
}

// Url parse link for special video ID
func Url(link *string) error {
	temp, err := url.Parse(*link)
	if err != nil {
		return err
	}
	if temp.Host == "www.youtube.com" {
		*link = temp.RawQuery[2:]
	} else if temp.Host == "youtu.be" {
		*link = temp.Path[1:]
	}
	return nil
}
