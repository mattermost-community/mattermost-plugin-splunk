package splunk

import (
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

func URLLastFragment(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", errors.Wrap(err, "wrong url")
	}
	idx := strings.LastIndex(u.Path, "/")
	return u.Path[idx+1:], nil
}
