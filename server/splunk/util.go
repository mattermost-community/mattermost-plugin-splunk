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

func findInSlice(ss []string, key string) int {
	for i, s := range ss {
		if s == key {
			return i
		}
	}
	return -1
}

func deleteFromSlice(ss []string, ind int) []string {
	if ind < 0 || ind >= len(ss) {
		return ss
	}

	ss[ind] = ss[len(ss)-1]
	ss[len(ss)-1] = ""
	return ss[:len(ss)-1]
}
