package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func getURLParam(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return "", errors.New("key missing")
	}

	return keys[0], nil
}
