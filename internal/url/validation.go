package url

import (
	"net/url"
)

func ValidateFilePath(path *string) error {

	_, err := url.ParseRequestURI(*path)
	if err != nil {
		return err
	}

	u, err := url.Parse(*path)
	if err != nil {
		return err
	}

	if err != nil || u.Scheme == "" || u.Host == "" {
		return err
	}

	// sanitize the url

	return nil

}
