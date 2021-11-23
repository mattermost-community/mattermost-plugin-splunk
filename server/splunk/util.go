package splunk

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func urlLastFragment(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", errors.Wrap(err, "wrong url")
	}
	idx := strings.LastIndex(u.Path, "/")
	return u.Path[idx+1:], nil
}
