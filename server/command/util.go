package command

import (
	"net/url"

	"github.com/pkg/errors"
)

func parseServerURL(u string) (string, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrap(err, "bad url")
	}
	if ur.Scheme != "http" && ur.Scheme != "https" {
		return "", errors.New("bad scheme")
	}
	ur.Scheme = "https"

	return ur.Scheme + "://" + ur.Host, err
}
