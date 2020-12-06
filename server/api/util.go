package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func getURLParam(r *http.Request, key string) (string, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return "", errors.New("key missing")
	}

	return val, nil
}
