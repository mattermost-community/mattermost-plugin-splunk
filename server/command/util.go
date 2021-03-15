package command

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func parseServerURL(u string, withPort bool) (string, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return "", errors.Wrap(err, "bad url")
	}
	if ur.Scheme != "http" && ur.Scheme != "https" {
		return "", errors.New("bad scheme")
	}
	ur.Scheme = "https"

	host := ur.Host
	if strings.Contains(host, ":") {
		i := strings.LastIndex(host, ":")
		host = host[:i]
	}
	if withPort {
		host += ":8089"
	}

	return ur.Scheme + "://" + host, nil
}
